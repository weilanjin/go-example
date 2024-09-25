package syncmap

import (
	"sync"
	"sync/atomic"
)

// 在下面两个场景中, sync.Map 会比 map + RWMutex 更高效:
// 1.在只会增长的缓存中, 一个key只被写入一次而被读多次.
// 2.多个goroutine为不相交(disjoint)的键集读、写和重写键值对.

// sync.Map 提供了9个方法
// 归为三类
// - 读操作
// 1. Load(key any) (value any, ok bool) // 读取一个键对应的值
// 2. Range(f func(key, value any) bool) // 遍历map, f返回false时停止遍历
// - 写操作
// 1. Store(key, value any) // 存储或者更新一个键值对
// 2. Delete(key any) // 删除一个键值对
// 3. Swap(key, value any) (previous any, loaded bool) // 交换键值对, 返回旧值, 如果这个键不存在, 返回false
// - 读写操作
// 1. CompareAndDelete(key, old any) (deleted bool) // 如果提供的键和旧值相等, 删除键值对
// 2. CompareAndSwap(key, old, new any) (swapped bool) // CAS 操作, 如果提供的键和旧值相等, 更新键值对
// 3. LoadOrStore(key, value any) (actual any, loaded bool) // 如果键存在, 返回键对应的值, 否则存储键值对
// 4. LoadAndDelete(key any) (value any, loaded bool) // 如果键存在, 返回键对应的值, 删除键值对

// sync.Map 的实现
// 1.[空间换时间] 一个dirty map, 一个read map 减少锁的粒度
// 2. read map read-only, 读操作时不需要加锁. 优先从 read 字段读取、更新、删除, 因为对 read 字段的读取不需要加锁.
// 3.[动态调整] 未命中次数多了之后, 将dirty数据提升为read数据, 避免总是从dirty中加锁读取.
// 4.[双重检查] 加锁之后,还要再检查read字段, 确定所查询的键值真的不存在, 才从dirty中读取.
// 5.[延迟删除] 删除一个键值只是打上标记,只有在创建dirty字段的时候才释放这个键, 然后, 只有在dirty数据提升为read数据的时候, read字段才会被删除.

type Map struct {
	mu sync.Mutex // 万不得已才使用的锁
	// 实际上是一个 “只读” 的map, 访问它的元素不需要加锁, 所以很快.
	read atomic.Pointer[readOnly]
	// 包含 map 中所有的元素, 包括新怎的元素
	// 访问 dirty 字段必须加锁, 当未命中达到一定的次数后会把它转为 read 字段
	dirty map[any]*entry
	// 未命中的次数表示有多少次数是未命中的, (不存在的元素)
	misses int
}

type readOnly struct {
	m       map[any]*entry
	amended bool // 如果有dirty数据,则返回true(有些数据只在dirty字段中,不在这个m中)
}

// expunged 标记一个键值对已经从 dirty 字段中删除了, 将它的值暂时设置为 expunged 这个值
// 标志出来
var expunged = new(any)

type entry struct{ p atomic.Pointer[any] } // 代表一个键值

func newEntry(i any) *entry {
	e := &entry{}
	e.p.Store(&i)
	return e
}

func (e *entry) load() (value any, ok bool) {
	p := e.p.Load()
	if p == nil || p == expunged {
		return nil, false
	}
	return *p, true
}

func (e *entry) trySwap(i *any) (*any, bool) {
	for {
		p := e.p.Load()
		// If the entry is expunged, trySwap returns false and leaves the entry
		// unchanged.
		if p == expunged {
			return nil, false
		}
		if e.p.CompareAndSwap(p, i) {
			return p, true
		}
	}
}

func (e *entry) unexpungeLocked() (wasExpunged bool) {
	return e.p.CompareAndSwap(expunged, nil)
}

func (e *entry) swapLocked(i *any) *any {
	return e.p.Swap(i)
}

func (e *entry) tryExpungeLocked() (isExpunged bool) {
	p := e.p.Load()
	for p == nil {
		if e.p.CompareAndSwap(nil, expunged) {
			return true
		}
		p = e.p.Load()
	}
	return p == expunged
}

func (m *Map) Swap(key, value any) (previous any, loaded bool) {
	read := m.loadReadOnly() // 1.读取 read 字段,原子操作, 没有使用到锁mu.
	// 2. 检查 read 字段是否包含这个key, 包含直接尝试交换
	if e, ok := read.m[key]; ok {
		if v, ok := e.trySwap(&value); ok { // 3. 如果key不存就是新增(虽然是读写操作但是不需要加锁)
			if v == nil {
				return nil, false
			}
			return *v, true
		}
	}
	m.mu.Lock()                   // 4. read 中不存在 开始加锁
	read = m.loadReadOnly()       // 5. 在临界区再次读取(double check)
	if e, ok := read.m[key]; ok { // 6. 如果存在,直接更新交换
		if e.unexpungeLocked() {
			m.dirty[key] = e
		}
		if v := e.swapLocked(&value); v != nil {
			m.dirty[key] = e
		}
	} else if e, ok := m.dirty[key]; ok { // 7. 处理存在 dirty 字段中key
		if v := e.swapLocked(&value); v != nil {
			loaded = true
			previous = *v
		}
	} else { // 8. read 和 dirty 都不存在的情况
		if !read.amended {
			m.dirtyLocked()
			m.read.Store(&readOnly{m: read.m, amended: true})
		}
		m.dirty[key] = newEntry(&value)
	}
	m.mu.Unlock()
	return
}

func (m *Map) Load(key any) (value any, ok bool) {
	read := m.loadReadOnly()
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read = m.loadReadOnly()
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = m.dirty[key]
			m.missLocked()
		}
		m.mu.Unlock()
	}
	if !ok {
		return nil, false
	}
	return e.load()
}

func (m *Map) missLocked() {
	m.misses++
	if m.misses < len(m.dirty) {
		return
	}
	m.read.Store(&readOnly{m: m.dirty})
	m.dirty = nil
	m.misses = 0
}

func (m *Map) dirtyLocked() {
	if m.dirty != nil {
		return
	}
	read := m.loadReadOnly()
	m.dirty = make(map[any]*entry, len(read.m))
	for k, e := range read.m {
		if !e.tryExpungeLocked() {
			m.dirty[k] = e
		}
	}
}

func (m *Map) loadReadOnly() readOnly {
	if p := m.read.Load(); p != nil {
		return *p
	}
	return readOnly{}
}
