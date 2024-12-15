package project

import (
	"sync"
	"sync/atomic"
)

// 双重检查 double check
// sync.Once
// sync.Map

type Once struct {
	m    sync.Mutex
	done uint32
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 0 { // 1
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 { // 2 拿到锁之后再次检查是否已经初始化了
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}