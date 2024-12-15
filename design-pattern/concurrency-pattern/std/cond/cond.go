package cond

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

/*
	type Cond
		func NewCond(l Locker) *Cond
		func (c *Cond) Broadcast()
		func (c *Cond) Signal()
		func (c *Cond) Wait()
*/

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

// 代码实现位于  runtime/sema.go 中
// 使用 平衡树 sudog 维护调用者列表
type notifyList struct {
	wait   uint32
	notify uint32
	lock   uintptr
	head   unsafe.Pointer
	tail   unsafe.Pointer
}

func runtime_notifyListAdd(l *notifyList) uint32
func runtime_notifyListWait(l *notifyList, t uint32)
func runtime_notifyListNotifyAll(l *notifyList)
func runtime_notifyListNotifyOne(l *notifyList)

type Cond struct {
	noCopy noCopy

	// 在检查条件或者修改条件时需要持有锁
	L       sync.Locker
	notify  notifyList
	checker copyChecker
}

func NewCond(l sync.Locker) *Cond {
	return &Cond{L: l}
}

type copyChecker uintptr

func (c *copyChecker) check() { // 检查没有被复制
	if uintptr(*c) != uintptr(unsafe.Pointer(c)) &&
		!atomic.CompareAndSwapUintptr((*uintptr)(c), 0, uintptr(unsafe.Pointer(c))) &&
		uintptr(*c) != uintptr(unsafe.Pointer(c)) { // 双重检查
		panic("sync.Cond is copied")
	}
}

// 允许调用者唤醒所有等待此 cond 和 goroutine

// 1.如果没有等待的goroutine, 则无须通知
// 2.如果有等待的goroutine, 则清空队列, 并唤醒所有等待的goroutine
// 主: 不强求调用者一定持有c.L的锁
func (c *Cond) Broadcast() { // notify all
	c.checker.check()
	runtime_notifyListNotifyAll(&c.notify) // 通知所有的waiter
}

// 允许调用者唤醒一个等待此Cond的goroutine
// 如果Cond的等待队列中有一个或者多个等待的goroutine, 则需要从等待队列中移除等待第一个goroutine, 并唤醒它
// 主: 不强求调用者一定持有c.L的锁
func (c *Cond) Signal() { // notify one
	c.checker.check()
	runtime_notifyListNotifyOne(&c.notify) // 通知一个 waiter
}

// 把调用者放入Cond的等待队列中并阻塞.直到被Signal或Broadcast方法从等待队列中移除并唤醒
// 主: 调用者必须持用c.L的锁
func (c *Cond) Wait() { // wait
	c.checker.check()
	t := runtime_notifyListAdd(&c.notify) // 先加入通知列表中
	c.L.Unlock()
	runtime_notifyListWait(&c.notify, t) // 等待通知
	c.L.Lock()                           // 唤醒后需要获取锁
}
