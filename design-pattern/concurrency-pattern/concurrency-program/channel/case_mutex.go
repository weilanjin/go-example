package channel

import "time"

// Mutex 使用channel实现互斥锁🔒
type Mutex struct {
	ch chan struct{}
}

func NewMutex() *Mutex {
	mu := &Mutex{ch: make(chan struct{}, 1)}
	mu.ch <- struct{}{} // 谁能取走,谁持有锁, 把值放回去就是释放锁
	return mu
}

// Lock 请求锁,直到获取到锁
func (m *Mutex) Lock() {
	<-m.ch
}

// UnLock 解锁
func (m *Mutex) UnLock() {
	select {
	case m.ch <- struct{}{}:
	default:
		panic("unlock if unlocked mutex")
	}
}

// TryLock 尝试获取锁
func (m *Mutex) TryLock() bool {
	select {
	case <-m.ch:
		return true
	default:
	}
	return false
}

// LockTimeout 获取锁,提供超时功能
func (m *Mutex) LockTimeout(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	select {
	case <-m.ch:
		timer.Stop()
		return true
	case <-timer.C:
	}
	return false
}

// IsLocked 锁是否已被持有
func (m *Mutex) IsLocked() bool {
	return len(m.ch) == 0
}