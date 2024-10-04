package ctx

import "context"

// 1.撤销动作一般都是主goroutine主动执行的
// 2.子goroutine需要主动检查上下文, 才能获知主goroutine是否下发了撤销命令

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
