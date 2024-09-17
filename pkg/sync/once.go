package sync

import (
	"sync"
	"sync/atomic"
)

// Do 返回可以 error

type Once struct {
	done uint32
	m    sync.Mutex
}

func (o *Once) Do(f func() error) error {
	if atomic.LoadUint32(&o.done) == 1 { // 快速路径
		return nil
	}
	return o.doSlow(f)
}

func (o *Once) doSlow(f func() error) error {
	o.m.Lock()
	defer o.m.Unlock()
	var err error
	if o.done == 0 { // 双重检查(double-checking), 获取到锁后检查是否同时已经有goroutine初始化了
		err = f() // 调用初始化函数
		if err == nil {
			atomic.StoreUint32(&o.done, 1) // 最后更改 done 值,表明已经初始化了
		}
	}
	return err
}

// Done 返回是否被执行过
// 如果被执行过,返回 true
// 如果没有执行过或者正在执行,则返回false
func (o *Once) Done() bool {
	// type Once struct {
	//   sync.Once
	// }
	// atomic.LoadUint32((*uint32)(unsafe.Pointer(&o.Once))) == 1 // 如果是内嵌标准库,可以这样取
	return atomic.LoadUint32(&o.done) == 1
}

func (o *Once) Reset() {
	atomic.StoreUint32(&o.done, 0)
}

// OnceFn 返回一个函数,该函数只会被执行一次
func OnceFn[T any](f func() T) func() T {
	var once sync.Once
	var t T
	fn := func() {
		t = f()
	}
	return func() T {
		once.Do(fn)
		return t
	}
}
