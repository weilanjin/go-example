package ctx

import (
	"context"
	"log"
	"time"
)

// Timeout 当前时间加上一段时间,最终会计算到未来的某个时间点
// Deadline 直接指明未来的某个时间点

func WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithDeadline(parent, time.Now().Add(timeout))
}

// Context 的 Done 方法返回的 channel 会在以下三种情况下关闭
// 1.当达到截止时间时
// 2.当返回的cancel函数被调用时
// 3.当父Context的Done channel被关闭时 (超时时间超过父节点的截止时间, 父节点时间到了就会直接取消)
func WithDeadline(parent context.Context, d time.Time) (context.Context, context.CancelFunc) {
	return nil, nil
}

func Usecase() {
	// case1: 超时
	log.Println("case1: expire")
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	<-ctx.Done()
	log.Println("err: ", ctx.Err())
	cancel()

	// case2: 主动取消
	log.Println("case2: cancel")
	ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	cancel()
	<-ctx.Done()
	log.Println("err: ", ctx.Err())

	// case3: 父context取消
	log.Println("case3: parent cancel")
	pCtx, pCancel := context.WithCancel(context.Background())
	ctx, cancel = context.WithDeadline(pCtx, time.Now().Add(5*time.Second))
	pCancel()
	<-ctx.Done()
	log.Println("err: ", ctx.Err())
	cancel()

	// case4: 超时超过父节点的截止时间
	log.Println("case4: parent expire")
	pCtx, pCancel = context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
	defer pCancel()
	ctx, cancel = context.WithDeadline(pCtx, time.Now().Add(5*time.Second))
	defer cancel()

	deadline, _ := ctx.Deadline()

	log.Println("timeout: ", time.Since(deadline))
}
