package context

import "time"

// Background() and TODO() 都是 emptyCtx{}

var (
	background = new(emptyCtx)
	todo       = new(emptyCtx)
)

func Background() Context {
	return background
}

func TODO() Context {
	return todo
}

// WithCancel
// 调用 cancel() 触发取消通知
func WithCancel(parent Context) (ctx Context, cancel func()) {
	c := cancelCtx{Context: parent}
	propagateCancel(parent, &c)
	return &c, func() {
		c.cancel(true, Canceled) // errors.New("context canceled")
	}
}

func propagateCancel(parent Context, child canceler) {
	done := parent.Done()
	select {
	case <-done:
		// parent is already canceled
		child.cancel(false, parent.Err())
		return
	default:
	}
	// p.children[child] = struct{}{}
}

func WithDeadline(parent Context, d time.Time) (Context, func()) {
	if cur, ok := parent.Deadline(); ok && cur.Before(d) /* 时间小于当前时间 */ {
		return WithCancel(parent)
	}
	c := &timerCtx{
		cancelCtx: cancelCtx{Context: parent},
		deadline:  d,
	}
	propagateCancel(parent, c)
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err == nil {
		// 定时器，到时间触发取消
		c.timer = time.AfterFunc(time.Until(d), func() {
			c.cancel(true, DeadlineExceeded) // context deadline exceeded
		})
	}
	return c, func() {
		// 手动取消
		c.cancel(true, Canceled) // errors.New("context canceled")
	}
}

// WithTimeout
// ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
// defer cancel() // 可以在超时之前释放资源
func WithTimeout(parent Context, timeout time.Duration) (Context, func()) {
	return WithDeadline(parent, time.Now().Add(timeout))
}

func WithValue(parent Context, key, val any) Context {
	return &valueCtx{parent, key, val}
}