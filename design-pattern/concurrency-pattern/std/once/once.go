package once

import (
	"sync"
	"sync/atomic"
)

type Once struct {
	done uint32
	m    sync.Mutex
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 0 { // 如果还没有初始化,则进入 doSlow, 否则直接返回
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 { // 双重检查(double-checking), 获取到锁后检查是否同时已经有goroutine初始化了
		defer atomic.StoreUint32(&o.done, 1) // 最后更改 done 值,表明已经初始化了
		f()                                  // 调用初始化函数
	}
}
