package singleflight

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
)

var errGoexit = errors.New("runtime.Goexit was called")

// 使用 Mutex 和 map 来实现
// Mutex 并发读写保护
// map 用来保存正在处理 in flight 对同一个key的请求

// call 代表一个正在执行的请求,或者已经执行完的请求
type call struct {
	wg sync.WaitGroup
	// val 这个字段代表处理完成的值, 在 waitGroup 完成之前只会写一次
	// 在 waitGroup 完成之后读取这个值
	val any
	err error

	// forgotten 指示当前call在处理时是否忘记这个key
	forgotten bool
	dups      int
	chans     []chan<- Result
}

type Result struct {
	Val    interface{}
	Err    error
	Shared bool
}

type panicError struct {
	value interface{}
	stack []byte
}

func (e panicError) Error() string {
	return fmt.Sprintf("%v\n\n%s", e.value, e.stack)
}

func (p *panicError) Unwrap() error {
	err, ok := p.value.(error)
	if !ok {
		return nil
	}

	return err
}

func newPanicError(v interface{}) error {
	stack := debug.Stack()

	// The first line of the stack trace is of the form "goroutine N [status]:"
	// but by the time the panic reaches Do the goroutine may no longer exist
	// and its status will have changed. Trim out the misleading line.
	if line := bytes.IndexByte(stack[:], '\n'); line >= 0 {
		stack = stack[line+1:]
	}
	return &panicError{value: v, stack: stack}
}

// Group 代表一个SingleFlight对象
type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	// 检查此key是否有执行中的任务
	if c, ok := g.m[key]; ok {
		c.dups++ // 重复任务 +1
		g.mu.Unlock()
		c.wg.Wait() // 等待正在执行的函数fn完成任务
		if e, ok := c.err.(*panicError); ok {
			panic(e)
		} else if errors.Is(c.err, errGoexit) {
			runtime.Goexit()
		}
		return c.val, c.err, true
	}
	c := new(call) // 没有执行中的任务, 它是第一个
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	g.doCall(c, key, fn) // 调用方法, 执行任务
	return c.val, c.err, c.dups > 0
}

func (g *Group) doCall(c *call, key string, fn func() (interface{}, error)) {
	normalReturn := false
	recovered := false
	defer func() {
		if !normalReturn && !recovered {
			c.err = errGoexit
		}
		g.mu.Lock()
		defer g.mu.Unlock()
		c.wg.Done()
		if g.m[key] == c { // 执行完毕,删除此key
			delete(g.m, key)
		}

		if e, ok := c.err.(*panicError); ok {
			if len(c.chans) > 0 {
				go panic(e)
				select {}
			} else {
				panic(e)
			}
		} else if errors.Is(c.err, errGoexit) {
		} else {
			// 正常返回,告诉那些waiter调用结果来了
			for _, ch := range c.chans {
				ch <- Result{
					Val:    c.val,
					Err:    c.err,
					Shared: c.dups > 0,
				}
			}
		}
	}()
	func() {
		defer func() {
			if !normalReturn {
				if r := recover(); r != nil {
					c.err = newPanicError(r)
				}
			}
		}()
		c.val, c.err = fn()
		normalReturn = true
	}()
	if !normalReturn {
		recovered = true
	}
}

func (g *Group) Forget(key string) {
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()
}