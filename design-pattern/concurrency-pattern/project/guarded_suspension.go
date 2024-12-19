package project

import (
	"fmt"
	"sync"
)

// 保护式挂起模式 Guarded Suspension
// 并发编程中使用的同步技术, 确保一个线程在继续执行前等待某个特定的条件变为真
//
// 在保护式挂起模式中, 一个线程在继续执行任务之前,会检查一个特定条件是否为真,如果条件为假,
// 那么线程会被挂起或检查被阻塞,直到条件变为真.

// net/rpc/client.go
func (client *Client) Call(serviceMethod string, args any, reply any) error {
	// 通过 Call Done 通知调用者,在通知之前,调用者一直处于保护挂起的状态
	call := <-client.Go(serviceMethod, args, reply, make(chan *Call, 1)).Done
	return call.Error
}

// Guard 方法提供了对函数fn的保护
// 如果没有获取到锁,调用者就会被阻塞挂起,一旦条件成熟,调用者获取到锁,并执行fn
func Guard(lock sync.Locker, fn func()) {
	lock.Lock()
	defer lock.Unlock()
	fn()
}

func GuardError(err *error) {
	if r := recover(); r != nil { // 使用 recover 方法捕获函数 panic
		if e, ok := r.(error); ok {
			*err = e
		} else {
			*err = fmt.Errorf("panic: %v", r)
		}
	}
}

func GuardFunc(fn func()) { // 避免函数fn发生panic,导致程序崩溃
	defer func() {
		recover()
	}()
	fn()
}

// GuardClose
// 关闭一个channel,防止panic
// 1.关闭一个值为nil的channel
// 2.关闭一个已经关闭的channel
func GuardClose[T any](ch chan T) { // 避免关闭已经关闭的channel
	defer func() {
		recover()
	}()
	close(ch)
}

func GuardSend[T any](ch chan<- T, val T) { // 避免发送到已经关闭的channel
	defer func() {
		recover()
	}()
	ch <- val
}