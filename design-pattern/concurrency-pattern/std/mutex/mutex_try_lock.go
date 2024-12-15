package mutex

import "sync/atomic"

// 2018 mutex 添加 TryLock 方法
// 2013 最早提出 #6123
// 2021 提交 #45435

// RWMutex TryLock 和 TryRLock 尝试写锁和读锁
// 返回 true 表示获取锁成功，false 表示获取锁失败
// 在已加锁的情况下, TryLock 将返回 false

func (m *MutexV120) TryLock() bool {
	old := m.state
	if old&(mutexLocked|mutexStarving) != 0 {
		return false
	}

	// 尝试设置加锁标志
	if !atomic.CompareAndSwapInt32(&m.state, old, old|mutexLocked) {
		return false
	}
	return true
}
