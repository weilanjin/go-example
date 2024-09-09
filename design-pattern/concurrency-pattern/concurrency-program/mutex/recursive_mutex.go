package mutex

import (
	"sync"
	"sync/atomic"

	"github.com/kortschak/goroutine"
)

// 递归锁, 也叫可重入锁
type RecursiveMutex struct {
	sync.Mutex
	owner     int64
	recursion int64
}

func (m *RecursiveMutex) Lock() {
	gid := goroutine.ID()
	if atomic.LoadInt64(&m.owner) == gid { // 如果当前加锁的goroutine就是此个goroutine
		atomic.AddInt64(&m.recursion, 1) // 递归/重入次数加1, 返回
		return
	}
	m.Mutex.Lock() // 尝试获取锁
	// 获取到锁,并且是第一次重入
	atomic.StoreInt64(&m.owner, gid)
	atomic.StoreInt64(&m.recursion, 1)
}

func (m *RecursiveMutex) Unlock() {
	gid := goroutine.ID()
	if atomic.LoadInt64(&m.owner) != gid { // 只允许加锁的goroutine解锁
		panic("unlock of unlocked mutex")
	}
	r := atomic.AddInt64(&m.recursion, -1) // 递归/重入次数减1
	if r != 0 {                            // 还需要递归释放锁
		return
	}
	// 释放锁
	atomic.StoreInt64(&m.owner, -1)
	m.Mutex.Unlock()
}
