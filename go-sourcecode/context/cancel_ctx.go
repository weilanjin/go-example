package context

import (
	"sync"
	"sync/atomic"
)

// 同时也会取消子 ctx
// 实现了 canceler 和 Context 接口
type cancelCtx struct {
	Context
	mu       sync.Mutex
	done     atomic.Value // chan struct{}
	children map[canceler]struct{}
	err      error
}

// itself
var cancelCtxKey int

func (c *cancelCtx) Value(key any) any {
	if key == &cancelCtxKey {
		return c
	}
	return value(c.Context, key)
}

func value(c Context, key any) any {
	for {
		switch ctx := c.(type) {
		// *valueCtx -> ctx.val
		case *valueCtx:
			if key == ctx.key {
				return ctx.val
			}
			c = ctx.Context
			// *cancelCtx -> c
			// *timerCtx -> &ctx.cancelCtx
			// *emptyCtx -> nil
			// default c.Value(key)
		}
	}
}

// Done return chan struct{}
func (c *cancelCtx) Done() <-chan struct{} {
	d := c.done.Load()
	if d != nil {
		return d.(chan struct{})
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	d = c.done.Load()
	if d == nil {
		d = make(chan struct{})
		c.done.Store(d)
	}
	return d.(chan struct{})
}

func (c *cancelCtx) Err() error {
	c.mu.Lock()
	err := c.err
	c.mu.Unlock()
	return err
}

// c.done
// 关闭每一个子节点
func (c *cancelCtx) cancel(removeFormParent bool, err error) {
	c.mu.Lock()
	if c.err != nil {
		c.mu.Unlock()
		return // 已经取消
	}
	c.err = err
	d, _ := c.done.Load().(chan struct{})
	if d == nil {
		c.done.Store(make(chan struct{}))
	} else {
		close(d)
	}
	for child := range c.children {
		child.cancel(false, err)
	}
	c.children = nil
	c.mu.Unlock()
	if removeFormParent {
		removeChild(c.Context, c)
	}
}

func removeChild(parent Context, child canceler) {
	p, ok := parentCancelCtx(parent)
	if !ok {
		return
	}
	p.mu.Lock()
	if p.children != nil {
		delete(p.children, child)
	}
	p.mu.Unlock()
}

func parentCancelCtx(parent Context) (*cancelCtx, bool) {
	done := parent.Done()
	if done == nil || done == make(chan struct{}) {
		return nil, false
	}
	p, ok := parent.Value(&cancelCtxKey).(*cancelCtx)
	if !ok {
		return nil, false
	}
	pdone, _ := p.done.Load().(chan struct{})
	if pdone != done {
		return nil, false
	}
	return p, true
}