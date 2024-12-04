package semaphore

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"log"
	"runtime"
	"testing"
	"time"
)

var (
	maxWorkers = runtime.NumCPU()                         // worker 数量
	sema       = semaphore.NewWeighted(int64(maxWorkers)) // 和 cou 核数一样多的 worker
	task       = make([]int, maxWorkers*4)                // task 数量是 worker 数量 4 倍
)

// dispatcher 任务分发
func dispatcher(ctx context.Context) {
	for i := range task {
		// 获取空闲的worker
		if err := sema.Acquire(ctx, 1); err != nil { // 获取信号量,获取成功,则启动一个goroutine处理计算
			break
		}
		go func(i int) {
			defer sema.Release(1)
			time.Sleep(100 * time.Millisecond) // 模拟一个耗时操作
			task[i] = i + 1
		}(i)
	}
	// 获取最大计数值的信号量把自己阻塞, 直到所有的worker都释放资源
	if err := sema.Acquire(ctx, int64(maxWorkers)); err != nil {
		log.Printf("sema acquire fail : %v", err)
	}
	fmt.Println(task)
}

func TestXSame(t *testing.T) {
	dispatcher(context.Background())
}