package ctx

import (
	"context"
	"time"
)

// DeadlineExceeded is the error returned by [Context.Err] when the context's
// deadline passes.
var DeadlineExceeded error = deadlineExceededError{}

type deadlineExceededError struct{}

func (deadlineExceededError) Error() string   { return "context deadline exceeded" }
func (deadlineExceededError) Timeout() bool   { return true }
func (deadlineExceededError) Temporary() bool { return true }

// 先检查父 Context的 截止时间, 如果此 Context 的 截止时间晚于父 Context 的截止时间, 则使用
// WithCancel(parent) 创建一个 Context 就好.
func WithDeadline(parent context.Context, d time.Time) (context.Context, context.CancelFunc) {
	if parent == nil {
		panic("cannot create a child of nil Context")
	}
	if cur, ok := parent.Deadline(); ok && cur.Before(d) {
		// 如果父Context的时间截止在这个时间d之前, 则应该使用父Context的截止时间
		return WithCancel(parent)
	}

	// 否则, 创建一个与时间相关的Context, 内部使用cancelCtx
	c := &timerCtx{
		cancelCtx: cancelCtx{Context: parent},
		deadline:  d,
	}

	propagateCancel(parent, c) // 向上传播
	dur := time.Until(d)       // time.Sub(time.Now())
	if dur <= 0 {              // 如果截止时间已过去
		c.cancel(true, DeadlineExceeded, nil) // 发生超时, 撤销这个Context
		return c, func() { c.cancel(false, Canceled, nil) }
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err == nil { // 设置一个定时器, 超过截止时间就撤销
		c.timer = time.AfterFunc(dur, func() {
			c.cancel(true, DeadlineExceeded, nil)
		})
	}
	return c, func() { c.cancel(true, Canceled, nil) }
}
