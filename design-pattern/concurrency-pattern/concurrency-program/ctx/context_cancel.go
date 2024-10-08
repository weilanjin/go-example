package ctx

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

func CancelCase() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-ctx.Done():
			return
		default:
			for {
				// 一段长时间运行,无法中途中止的代码
			}
		}
	}()
	cancel()
}

// 1.撤销动作一般都是主goroutine主动执行的
// 2.子goroutine需要主动检查上下文, 才能获知主goroutine是否下发了撤销命令

type canceler interface {
	cancel(removeFromParent bool, err, cause error)
	Done() <-chan struct{}
}

var Canceled = errors.New("context canceled")
var cancelCtxKey int

// closedchan is a reusable closed channel.
var closedchan = make(chan struct{})

func init() {
	close(closedchan)
}

type cancelCtx struct {
	Context

	mu       sync.Mutex
	done     atomic.Value
	children map[canceler]struct{} // 此 context 的子 Context 对象
	err      error                 // 撤销时设置的error
	cause    error                 // 撤销时的根因
}

func newCancelCtx(parent context.Context) *cancelCtx {
	return &cancelCtx{
		Context: parent,
	}
}

func (c *cancelCtx) cancel(removeFromParent bool, err, cause error) {
	if err == nil {
		panic("context: internal error: missing cancel error")
	}
	if cause == nil {
		cause = err
	}
	c.mu.Lock()
	if c.err != nil {
		c.mu.Unlock()
		return // 已经被撤销过,直接返回即可
	}
	c.err = err
	c.cause = cause
	d, _ := c.done.Load().(chan struct{}) // 读取 done 这个 channel
	if d == nil {
		c.done.Store(closedchan) // 既然已经明确被撤销了, 那么直接使用一个已关闭的channel即可
	} else {
		close(d)
	}
	for child := range c.children { // 子 Context 也都要被撤销
		child.cancel(false, err, cause)
	}
	c.children = nil
	c.mu.Unlock()

	if removeFromParent {
		removeChild(c.Context, c) // 清空子 Context 列表
	}
}

func removeChild(parent context.Context, child canceler) {
	if s, ok := parent.(stopCtx); ok {
		s.stop()
		return
	}
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

func (c *cancelCtx) Done() <-chan struct{} {
	d := c.done.Load()
	if d != nil { // 如果已经初始化done, 直接返回
		return d.(chan struct{})
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	d = c.done.Load()
	if d == nil { // 双重检查, 如果还没有初始化过,则初始化一个未关闭的 channel
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

func (c *cancelCtx) Value(key any) any {
	if key == &cancelCtxKey { // 如果查询是自己
		return c
	}
	return value(c.Context, key) // 否则往上查找
}

func WithCancel(parent context.Context) (ctx context.Context, cancel context.CancelFunc) {
	c := withCancel(parent)
	return c, func() {
		c.cancel(true, Canceled, nil)
	}
}

func withCancel(parent context.Context) *cancelCtx {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	c := newCancelCtx(parent)
	propagateCancel(parent, c) // 向上传播,让父 Context 关联这个子 context
	return c
}

func propagateCancel(parent context.Context, child canceler) {
	done := parent.Done()
	if done == nil {
		return // 如果父 Context 永远不会被撤销, 比如 context.Background() 和 context.TODO(), 则不需要处理, 返回
	}

	select {
	case <-done:
		// 父 Context 已经被撤销, 这个子Context也要被撤销
		child.cancel(false, parent.Err(), Cause(parent))
		return
	default:
	}

	// 得到父 Context 的可撤销对象,或者往上找, 直到找到一个可撤销的Context, 或者不存在
	if p, ok := parentCancelCtx(parent); ok {
		p.mu.Lock()
		if p.err == nil {
			// 如果是父Context 已经被撤销了, 则当前这个Context也要被撤销
			child.cancel(false, p.err, p.cause)
		} else { // 否则, 把自己加入父 Context 的子 Context 列表中
			if p.children == nil {
				p.children = make(map[canceler]struct{})
			}
			p.children[child] = struct{}{}
		}
		p.mu.Unlock()
	} else { // 如果父 Context 以上都不是可能撤销的Context, 那么此 Context 自己启动一个 goroutine 监听
		// goroutines.Add(1)
		go func() {
			select {
			case <-parent.Done():
				child.cancel(false, parent.Err(), Cause(parent))
			case <-child.Done():
			}
		}()
	}
}

func parentCancelCtx(parent context.Context) (*cancelCtx, bool) {
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

func Cause(c context.Context) error {
	if cc, ok := c.Value(&cancelCtxKey).(*cancelCtx); ok {
		cc.mu.Lock()
		defer cc.mu.Unlock()
		return cc.cause
	}
	return c.Err()
}
