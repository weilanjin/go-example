package mutex

import (
	"sync/atomic"
)

// v 2008

// CAS 操作, 当时还没有抽象出 atomic 包
func cas(val *int32, old, new int32) bool { return false }

// semacquire 用来把调用者 goroutine 压入一个队列,并把此 goroutine 设置为阻塞状态
// 主要用来处理不能获取到锁的 goroutine (waiter, 等待者)
func semacquire(*int32) {}

// semrelease 用来从队列中取出一个 goroutine, 唤醒它, 并获得到锁.
func semrelease(*int32) {}

// 互斥锁的结构

type Mutex struct {
	/*
		key 持有锁和等待锁的数量
		0: 锁未被持有
		1: 锁被持有,没有等待者
		n: 锁被持有, 还有 n-1 个等待者

		sema
			等待者队列使用的信号量
	*/

	key  int32 // 锁是否被持有的标志 key >= 1
	sema int32 // 信号量专用, 用于阻塞/唤醒 goroutine
}

// 保证成功在 val 上增加delta的值
func xadd(val *int32, delta int32) (new int32) {
	for {
		v := *val
		if cas(val, v, v+delta) {
			return v + delta
		}
	}
	panic("unreached")
}

// 获取锁
func (m *Mutex) Lock() {
	if xadd(&m.key, 1) == 1 { // 1.将标志量加1, 如果等于1, 则表示称获取锁
		return
	}
	semacquire(&m.sema) // 2. 否则阻塞等待
}

func (m *Mutex) Unlock() {
	if xadd(&m.key, -1) == 0 { // 3.将标志量减1, 如果等于0,则表示没有其他waiter
		return
	}
	semrelease(&m.sema) // 4.唤醒其他被阻塞的 goroutine
}

// v 2008+ 微调
// 1. int32 -> uint32
// 2. 调用 Unlock 方法会做检查
// 3. 使用 atomic包同步原语执行原子操作
type MutexV2 struct {
	key  uint32
	sema uint32
}

// v 2011
// v2 goroutine 会排队等待获取互斥锁. 性能不是最优
// 1. 改进新来的goroutine也可能获取到锁的机会 (但是失败后加入等待队列中)
type Mutex2011 struct {
	state int32 // mutexWaiters | mutexWoken | mutexLocked
	sema  uint32
}

const (
// mutexLocked      = 1 << iota // 第一位代表是否加锁  1 001
// mutexWoken                   // 唤醒标志          2  010
// mutexWaiterShift = iota      // waiter 开始的位   2  010
)

/*
+-----------------+------------------------------------------------------+---------------------+
| 请求锁的goroutine类型 | 已加锁                                                  | 未加锁                 |
+-----------------+------------------------------------------------------+---------------------+
| 新来的goroutine    | waiter++; 休眠等待                                       | 获取到锁                |
| 被唤醒的goroutine   | 新来的goroutine已经抢到锁,waiter++;清除mutexWoken标志;重新休眠,回到队列中 | 清除mutexWoken标志;获取到锁 |
+-----------------+------------------------------------------------------+---------------------+
*/
func (m *Mutex2011) Lock() {
	// 1.快速路径: 幸运,能够直接获取到锁(没有那个goruntine持有锁, 也没有等待持用锁的goroutine)
	if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
		return
	}
	awoke := false // 是否被唤醒
	for {          // 2. m.state 不是零值, 循环进行检查
		old := m.state
		new := old | mutexLocked  // 3.新状态已加锁
		if old&mutexLocked != 0 { // 4.锁已被持有, 加入waiter中
			new = old + 1<<mutexWaiterShift // 5.waiter 数量加1
		}
		if awoke { // 6.
			// 此goroutine是被唤醒的
			// 新的状态, 清除唤醒标志
			new &^= mutexWoken
		}
		if atomic.CompareAndSwapInt32(&m.state, old, new) { // 7.设置为最新状态
			if old&mutexLocked == 0 { // 8.锁的原状态, 未加锁, 在此 goroutine 获取了锁, 成功!
				break
			}
			// runtime.Semacquire(&m.sema) // 9.请求信号量, 加入队列中
			awoke = true // 10.被唤醒, 和新的 goroutine 抢锁
		}
	}
}

func (m *Mutex2011) Unlock() {
	// 快速路径: 去除锁的标志位
	new := atomic.AddInt32(&m.state, -mutexLocked) // 去除锁标志位 1.尝试将持有锁的标志设置为未加锁的状态
	if (new+mutexLocked)&mutexLocked == 0 {        // 2.本来就没有加锁
		panic("sync: unlock of unlocked mutex")
	}
	old := new
	for {
		if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken) != 0 {
			// 3.没有 waiter, 或者有被唤醒的waiter,或者原来已加锁
			break
		}
		new = (old - 1<<mutexWaiterShift) | mutexWoken // 新状态,准备唤醒goroutine,并设置唤醒标志
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			// runtime.Semrelease(&m.sema)
			break
		}
		old = m.state
	}
}

// v2015
// 自旋: 多几次尝试机会
// 1. 问题(饥饿)新来的 goroutine 也参与竞争,有可能每次都被新来的goroutine 抢到锁
type Mutex2015 struct {
	state int32
}

func runtime_canSpin(i int) bool                                   { return false }
func runtime_doSpin()                                              {}
func runtime_Semacquire(sema *int32)                               {}
func runtime_Semrelease(s *uint32, handoff bool, skipframes int)   {}
func runtime_nanotime() int64                                      { return 0 }
func runtime_SemacquireMutex(s *uint32, lifo bool, skipframes int) {}
func throw(s string)                                               {}
func fatal(s string)                                               {}

func (m *Mutex2015) Lock() {
	// 快捷路径 幸运之路, 正好获取到锁
	if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
		return
	}
	awoke := false
	iter := 0
	// 临界区的代码执行时间很短, 锁很快就能被释放,而请求锁的goroutine不用通过休眠唤醒的方式等待调度, 直接自旋几次,可能
	// 就获得到了锁,减少了goroutine上下文切换的开销.
	for {
		old := m.state            // 先保存当前锁的状态
		new := old | mutexLocked  // 新状态: 设置已加锁标志
		if old&mutexLocked != 0 { // 锁还没有被释放
			if runtime_canSpin(iter) { // 还可以自旋
				if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 &&
					atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
					awoke = true
				}
				runtime_doSpin()
				iter++
				continue // 自选, 再次尝试请求锁
			}
			new = old + 1<<mutexWaiterShift
		}
		if awoke { // 被唤醒
			if new&mutexWoken == 0 {
				panic("sync: inconsistent mutex state")
			}
			new &^= mutexWoken // 新状态, 清除唤醒标志
		}
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			if old&mutexLocked == 0 { // 旧状态, 锁已被释放; 新状态,成功持有了锁,直接返回
				break
			}
			runtime_Semacquire(&m.state) // 到一定次数还没有拿到锁,再去阻塞等待
			awoke = true                 // 被唤醒
			iter = 0
		}
	}
}

// v go 1.20
// 处理饥饿问题, 快速路径和慢速路径拆成独立的方法, 便以内联,提高性能
type MutexV120 struct {
	state int32 // waiter count、唤醒标志、锁状态
	sema  uint32
}

const (
	mutexLocked      = 1 << iota // 1 Mutex 加锁标志
	mutexWoken                   // 2
	mutexStarving                // 4
	mutexWaiterShift = iota      // 3

	// 1ms
	starvationThresholdNs = 1e6
)

func (m *MutexV120) Lock() {
	// 快速路径 获取到锁
	if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
		return
	}
	m.lockSlow()
}

// It is a run-time error if m is not locked on entry to Unlock.
func (m *MutexV120) Unlock() {
	// Fast path: drop lock bit.
	new := atomic.AddInt32(&m.state, -mutexLocked)
	if new != 0 {
		// Outlined slow path to allow inlining the fast path.
		// To hide unlockSlow during tracing we skip one extra frame when tracing GoUnblock.
		// 如果还有 waiter, 或者Mutex处于饥饿状态,则调用unlockSlow
		m.unlockSlow(new)
	}
}

func (m *MutexV120) lockSlow() {
	var (
		waitStartTime int64
		starving      bool // 是否处于饥饿
		awoke         bool
		iter          int
		old           = m.state
	)
	for {
		// Don't spin in starvation mode, ownership is handed off to waiters
		// so we won't be able to acquire the mutex anyway.
		// 在饥饿模式下不要自旋，因为锁的所有权已移交给等待者，
		// 所以我们无法获取该互斥锁。这意在避免无谓的CPU占用。
		if old&(mutexLocked|mutexStarving) == mutexLocked && runtime_canSpin(iter) {
			// Active spinning makes sense.
			// Try to set mutexWoken flag to inform Unlock
			// to not wake other blocked goroutines.
			// 尝试设置 mutexWoken 标志, 并通知Unlock不要唤醒其他阻塞的goroutine。
			if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 &&
				atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
				awoke = true
			}
			runtime_doSpin()
			iter++
			old = m.state
			continue
		}
		new := old
		// Don't try to acquire starving mutex, new arriving goroutines must queue.
		// 只有在非饥饿状态下才尝试获取锁; 否则,新的 goroutine 应该被加入 waiter 队列中
		if old&mutexStarving == 0 {
			new |= mutexStarving
		}
		if old&(mutexLocked|mutexStarving) != 0 {
			new += 1 << mutexWaiterShift
		}

		// The current goroutine switches mutex to starvation mode.
		// But if the mutex is currently unlocked, don't do the switch.
		// Unlock expects that starving mutex has waiters, which will not
		// be true in this case.
		// 当前 goroutine 尝试将锁设置为饥饿状态, 只有在当前的状态为已加锁的情况下才这么做
		if starving && old&mutexLocked != 0 {
			new |= mutexStarving
		}
		if awoke { // 清除唤醒标志
			if new&mutexWoken == 0 {
				throw("sync: inconsistent mutex state")
			}
			new &^= mutexWoken
		}

		// 设置新的状态
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			if old&(mutexLocked|mutexStarving) == 0 {
				break // locked the mutex with CAS 在非饥饿状态下获取到锁
			}
			// 以下处理饥饿状态

			// 如果 waiter 以前就在队列里面, 则使用 LIFO 策略
			queueLifo := waitStartTime != 0
			if waitStartTime == 0 {
				waitStartTime = runtime_nanotime()
			}
			runtime_SemacquireMutex(&m.sema, queueLifo, 1)                                  // 阻塞等待
			starving = starving || runtime_nanotime()-waitStartTime > starvationThresholdNs // 唤醒后检查是否饥饿
			old = m.state
			if old&mutexStarving != 0 {
				// If this goroutine was woken and mutex is in starvation mode,
				// ownership was handed off to us but mutex is in somewhat
				// inconsistent state: mutexLocked is not set and we are still
				// accounted as waiter. Fix that.
				// 非正常状态
				if old&(mutexLocked|mutexWoken) != 0 || old>>mutexWaiterShift == 0 {
					throw("sync: inconsistent mutex state")
				}

				delta := int32(mutexLocked - 1<<mutexWaiterShift)
				if !starving || old>>mutexWaiterShift == 1 { // 退出饥饿状态
					// Exit starvation mode.
					// Critical to do it here and consider wait time.
					// Starvation mode is so inefficient, that two goroutines
					// can go lock-step infinitely once they switch mutex
					// to starvation mode.
					delta -= mutexStarving
				}
				atomic.AddInt32(&m.state, delta)
				break
			}
			awoke = true
			iter = 0
		} else {
			old = m.state
		}
	}
}

func (m *MutexV120) unlockSlow(new int32) {
	if (new+mutexLocked)&mutexLocked == 0 {
		fatal("sync: unlock of unlocked mutex")
	}
	// 非饥饿状态
	if new&mutexStarving != 0 {
		old := new
		for {
			// If there are no waiters or a goroutine has already
			// been woken or grabbed the lock, no need to wake anyone.
			// In starvation mode ownership is directly handed off from unlocking
			// goroutine to the next waiter. We are not part of this chain,
			// since we did not observe mutexStarving when we unlocked the mutex above.
			// So get off the way.
			// 如果没有waiter, 或者已经有goroutine被唤醒了, 或者已经有goroutine抢到了锁,或者mutex处于饥饿状态,
			// 则不需要唤醒任何goroutine
			if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken|mutexStarving) != 0 {
				return
			}
			// Grab the right to wake someone.
			// 唤醒一个goroutine,并设置唤醒标志
			new = (old - 1<<mutexWaiterShift) | mutexLocked
			if atomic.CompareAndSwapInt32(&m.state, old, new) {
				runtime_Semrelease(&m.sema, false, 1)
				return
			}
			old = m.state
		}
	} else {
		// Starving mode: handoff mutex ownership to the next waiter, and yield
		// our time slice so that the next waiter can start to run immediately.
		// Note: mutexLocked is not set, the waiter will set it after wakeup.
		// But mutex is still considered locked if mutexStarving is set,
		// so new coming goroutines won't acquire it.
		// 启动饥饿模式: 将互斥锁的所有权交给下一个等待者,并让出当前的时间片.
		runtime_Semrelease(&m.sema, true, 1)
	}
}
