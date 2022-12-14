package context

import (
	"sync"
	"sync/atomic"
)

// 同时也会取消子 ctx
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

func (c *cancelCtx) cancel(removeFormParent bool, err error) {

}