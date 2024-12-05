// Package waitgroup
// 解决并发-等待问题
// 等所有goroutine执行完毕后才继续执行
// 例如:
// Linux - barrier (屏障)
// C++ - std::barrier
// Java - CyclicBarrier 和 CountDownLatch
// ---
// https://github.com/sourcegraph/conc
//
// 第三方库
package waitgroup

/*
	第三方库
	golang.org/x/sync/errgroup
		type Group
			func WithContext(ctx context.Context) (*Group, context.Context)
			func (g *Group) Go(f func() error)
			func (g *Group) SetLimit(n int)
			func (g *Group) TryGo(f func() error) bool
			func (g *Group) Wait() error

	github.com/mdlayher/schedgroup
		type Group
			func New(ctx context.Context) *Group
			func (g *Group) Delay(delay time.Duration, fn func())
			func (g *Group) Schedule(when time.Time, fn func())
			func (g *Group) Wait() error

	github.com/vardius/gollback
			func All(ctx context.Context, fns ...AsyncFunc) ([]interface{}, []error)
			func Race(ctx context.Context, fns ...AsyncFunc) (interface{}, error)
			func Retry(ctx context.Context, retires int, fn AsyncFunc) (interface{}, error)
			type AsyncFunc

	github.com/aaronjan/hunch
			func All(parentCtx context.Context, execs ...Executable) ([]interface{}, error)
			func Last(parentCtx context.Context, num int, execs ...Executable) ([]interface{}, error)
			func Retry(parentCtx context.Context, retries int, fn Executable) (interface{}, error)
			func Take(parentCtx context.Context, num int, execs ...Executable) ([]interface{}, error)
			func Waterfall(parentCtx context.Context, execs ...ExecutableInSequence) (interface{}, error)
			type Executable
			type ExecutableInSequence
			type IndexedExecutableOutput
			type IndexedValue
			type MaxRetriesExceededError
			func (err MaxRetriesExceededError) Error() string
*/