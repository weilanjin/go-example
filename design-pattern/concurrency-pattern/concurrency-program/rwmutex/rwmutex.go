package rwmutex

import (
	"sync"
	"sync/atomic"
)

type RWMutex struct {
	w           sync.Mutex   // 由 pending writer 持有这个锁; 当 writer 持有这个锁时, 它会持有这个互斥锁 w
	writerSem   uint32       // 为 writer 设置的信号量, writer 等待先前的 reader 释放锁; 用来阻塞以及唤醒 writer 的信号
	readerSem   uint32       // 为 reader 设置的信号量, reader 等待先前的 writer 释放锁; 用来阻塞以及唤醒 reader 的信号
	readerCount atomic.Int32 // pending reader 的数量; 当前 reader 的数量(包括持有读锁的 和 等待读锁的)
	readerWait  atomic.Int32 // departing reader 的数量; 当前持有读锁的数量
}

func runtime_SemacquireRWMutex(s *uint32, handoff bool, skipframes int) {}
func runtime_Semrelease(s *uint32, handoff bool, skipframes int)        {}

func (rw *RWMutex) RLock() {
	if rw.readerCount.Add(1) < 0 {
		// pending writer 等待中
		runtime_SemacquireRWMutex(&rw.readerSem, false, 0)
	}
}

func (rw *RWMutex) RUnlock() {
	if r := rw.readerCount.Add(-1); r < 0 {
		rw.rUnlockSlow(r)
	}
}

func (rw *RWMutex) rUnlockSlow(r int32) {
	// 如果有 pending 状态的writer
	if rw.readerWait.Add(-1) == 0 {
		// 最后一个 reader 唤醒writer
		runtime_Semrelease(&rw.writerSem, false, 1)
	}
}

const rwmutexMaxReaders = 1 << 30 // 1073741824

func (rw *RWMutex) Lock() {
	rw.w.Lock()
	r := rw.readerCount.Add(-rwmutexMaxReaders) + rwmutexMaxReaders
	// Wait for active readers.
	// 注意: 这一行“一箭双雕”,即把reader的数量变为负值, 又获得先前reader的数量
	if r != 0 && rw.readerWait.Add(r) != 0 {
		// 如果还有已经获取到读锁的reader, 那么这个writer就需要等待
		runtime_SemacquireRWMutex(&rw.writerSem, false, 0)
	}
}

func (rw *RWMutex) Unlock() {
	r := rw.readerCount.Add(rwmutexMaxReaders) // 把reader的数量变为正值
	for i := 0; i < int(r); i++ {
		runtime_Semrelease(&rw.readerSem, false, 0) // 唤醒那些等待释放写锁的reader, 解放它
	}
	rw.w.Unlock()
}

func (rw *RWMutex) TryRLock() bool {
	for {
		c := rw.readerCount.Load() // 获取当前reader的数量
		if c < 0 {
			return false // 当前有 writer 持有写锁, reader 不能获取读锁, 直接返回
		}
		if rw.readerCount.CompareAndSwap(c, c+1) { // reader 的数量加1, 获取到读锁
			return true
		}
	}
}
