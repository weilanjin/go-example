package rwmutex

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type RWMutexEx struct {
	sync.RWMutex
}

type m struct {
	w           sync.Mutex
	writerSem   uint32
	readerSem   uint32
	readerCount atomic.Int32
	readerWait  atomic.Int32
}

const (
	mutexLocked      = 1 << iota // 加锁标志位 1
	mutexWoken                   // 2
	mutexStarving                // 饥饿标志位 4
	mutexWaiterShift = iota      // 3
)

func (rw *RWMutexEx) ReaderCount() int {
	v := (*m)(unsafe.Pointer(&rw.RWMutex))
	r := v.readerCount.Load()
	if r < 0 {
		r += rwmutexMaxReaders
	}
	return int(r)
}

func (rw *RWMutexEx) ReaderWait() int {
	v := (*m)(unsafe.Pointer(&rw.RWMutex))
	c := v.readerWait.Load()
	return int(c)
}

func (rw *RWMutexEx) WriterCount() int {
	v := atomic.LoadInt32((*int32)(unsafe.Pointer(&rw.RWMutex)))
	v = v >> mutexWaiterShift
	v = v + (v & mutexLocked)
	return int(v)
}
