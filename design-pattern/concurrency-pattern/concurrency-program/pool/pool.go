package pool

import (
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

// Go 1.13 之前 sync.Pool 的实现,有两大问题
// 1. 每次垃圾回收时都会回收创建的对象
// - 如果缓存的对象数量太多,就会导致STW的时间变长
// - 缓存的对象都被回收后,则会导致Get命中率下降, Get 的方法不得重新创建很多对象.
// 2. 底层实现使用了Mutex,对这个锁并发请求竞争激烈的时候, 会导致性能下降
// Go 1.13
// 对 Pool 优化就是避免使用锁, 同时将加锁的队列改成 lock-free 的队列
// 给将移除的对象再对一次“复活”的机会.

type Pool struct {
	// 每次垃圾回收的时候, Pool都会把victim中的对象移除, 然后把local的数据给移到victim,
	// 这样一来, local就会被清空, 而victim就是一个垃圾分拣站,以后它里面的东西可能被当作垃
	// 圾丢弃,但是里面有用的东西也可能会被捡回来重新使用.

	// 当前所有的空闲对象都被放在local字段中, 在请求对象时, 也是优先从local字段中查找可用对象
	local     unsafe.Pointer
	localSize uintptr

	victim     unsafe.Pointer // local from previous cycle
	victimSize uintptr        // size of victims array

	New func() any
}

func (p *Pool) Put(x any) {
	if x == nil { // 首先检查放入的对象是否为nil,如果是nil就直接返回,因为放入nil没什么意义.
		return
	}
	l, _ := p.pin() // 把goroutine和P绑定,避免其他的 goroutine把该P抢去
	if l.private == nil {
		// 一个快速模式,如果private为nil,则直接赋值给它.
		// private 无须像本地local一样需要加锁或者local-free,可以更高效.
		l.private = x
	} else {
		// 如果private不为nil,则将对象压入本地的共享队列中.
		// 共享队列是一个lock-free 双向队列.
		l.shared.pushHead(x)
	}
	runtime_procUnpin()
}

func (p *Pool) Get() any {
	// 把当前的 goroutine 固定在当前的P上,这样一来, 在操作与这个P相关的对象时就不用加锁
	// 因为每个P都只有一个活动的goroutine在运行.
	// 每个P都有自动的缓存,优先从这个缓存中读/写对像,不要加锁如果没哟再去其他P中获取对象.
	// timer、goroutine 任务调度都是采用相同的原理
	l, pid := p.pin()
	x := l.private // 检查此P的local的private字段,如果存在,就使用这个对象.
	l.private = nil
	if x == nil { // 在 local.private 为空的情况下,检查本地的其他缓存队列.
		x, _ = l.shared.popHead() // 如果本地队列中有缓存的 对象,则返回改对象.
		if x == nil {
			// Get 方法首先尝试从其他的P中获取对象, 如果获取失败, 则从victim中
			x = p.getSlow(pid)
		}
	}
	runtime_procPin()
	// “复活”一个对象, 如果不成功,就创建一个新的对象.
	if x == nil && p.New != nil {
		x = p.New()
	}
	return x
}

func (p *Pool) getSlow(pid int) any {
	size := runtime_LoadAcquintptr(&p.localSize) // load-acquire
	locals := p.local                            // load-consume
	// 尝试从proc获取一个对象
	for i := 0; i < int(size); i++ { // 从其他的P中获取对象, 从下一个P开始依次检查,看看有没有缓存对象,如果有,则返回改对象.
		l := indexLocal(locals, (pid+i+1)%int(size))
		if x, _ := l.shared.popHead(); x != nil {
			return x
		}
	}
	// 尝试从victim 中获取对象
	// 检查victim和检查local的方式一样,毕竟它们时相同的类型,
	size = atomic.LoadUintptr(&p.victimSize)
	if uintptr(pid) >= size {
		return nil
	}
	locals = p.victim
	l := indexLocal(locals, pid)
	if x := l.private; x != nil {
		l.private = nil
		return x
	}
	for i := 0; i < int(size); i++ {
		l := indexLocal(locals, (pid+i)%int(size))
		if x, _ := l.shared.popTail(); x != nil {
			return x
		}
	}
	atomic.StoreUintptr(&p.victimSize, 0)
	return nil
}

func (p *Pool) pin() (*poolLocal, int) {
	if p == nil {
		panic("nil Pool")
	}

	pid := runtime_procPin()
	s := runtime_LoadAcquintptr(&p.localSize)
	l := p.local
	if uintptr(pid) < s {
		return indexLocal(l, pid), pid
	}
	return p.pinSlow()
}

func (p *Pool) pinSlow() (*poolLocal, int) {
	runtime_procUnpin()
	allPoolsMu.Lock()
	defer allPoolsMu.Unlock()
	pid := runtime_procPin()
	s := p.localSize
	l := p.local
	if uintptr(pid) < s {
		return indexLocal(l, pid), pid
	}
	if p.local == nil {
		allPools = append(allPools, p)
	}
	size := runtime.GOMAXPROCS(0)
	local := make([]poolLocal, size)
	atomic.StorePointer(&p.local, unsafe.Pointer(&local[0]))
	runtime_StoreReluintptr(&p.localSize, uintptr(size))
	return &local[pid], pid
}

func indexLocal(l unsafe.Pointer, i int) *poolLocal {
	lp := unsafe.Pointer(uintptr(l) + uintptr(i)*unsafe.Sizeof(poolLocal{}))
	return (*poolLocal)(lp)
}

type poolLocal struct {
	poolLocalInternal
	pad [128 - unsafe.Sizeof(poolLocalInternal{})%128]byte
}

type poolLocalInternal struct {
	// Can be used only by the respective P
	// 代表一个缓存对象
	// 只能由相应的那个P获取, 因为一个P同时只能执行一个goroutine
	// 不会有并发问题
	private any
	// local P can pushHead/popHead; any P can popTail
	// 可以任意的P访问
	shared poolChain
}

func runtime_StoreReluintptr(ptr *uintptr, val uintptr) {}

func runtime_procUnpin()
func runtime_LoadAcquintptr(ptr *uintptr) uintptr
func runtime_procPin() int { return 0 }
func runtime_registerPoolCleanup(cleanup func())

var (
	allPoolsMu sync.Mutex
	// allPools 是一组Pool对象,它们拥有非空的主缓存(non-empty primary cache).
	// 可以由allPoolsMu / pinning 或者 STW 保证并发安全
	allPools []*Pool
	// oldPool 是一组Pool对象, 它们拥有非空的victim缓存(non-empty victim cache)
	// 可以由STW 保证并发安全
	oldPools []*Pool
)

func init() {
	runtime_registerPoolCleanup(poolCleanup)
}

func poolCleanup() {
	// Dorp victim caches from all pools
	// 丢弃当前的 victim, 所以STW不用加锁
	for _, p := range oldPools {
		p.victim = nil
		p.victimSize = 0
	}
	// Move primary cache to victim cache
	// 把当前的 local 移动到 victim, 并将local清空
	for _, p := range allPools {
		p.victim = p.local
		p.victimSize = p.localSize
		p.local = nil
		p.localSize = 0
	}
	oldPools, allPools = allPools, nil
}
