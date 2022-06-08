package context

import (
	"context"
	"log"
	"testing"
	"time"
)

// 假设有一个超长的切片，切片的元素类型为int，切片中的元素为乱序排列。
// 限时5秒，使用多个goroutine查找切片中是否存在给定值，
// 在找到目标值或者超时后立刻结束所有goroutine的执行。
// 比如切片为：[23, 32, 78, 43, 76, 65, 345, 762, …… 915, 86]，
// 查找的目标值为345，如果切片中存在目标值程序输出:"Found it!"
// 并且立即取消仍在执行查找任务的goroutine。如果在超时时间未找到目标值程序输出:
// "Timeout! Not Found"，同时立即取消仍在执行查找任务的goroutine

func Test1(t *testing.T) {
	data := []int{23, 32, 78, 43, 76, 65, 345, 762, 915, 86}
	dataLen := len(data)
	size, target := 3, 345

	ctx, cancel := context.WithCancel(context.Background())
	resCh := make(chan int, 1)

	for i := 0; i < dataLen; i += size {
		end := i + size
		if end >= dataLen {
			end = dataLen - 1
		}
		go SearchTarget(ctx, data[i:end], target, resCh)
	}

	timer := time.NewTimer(time.Second * 2)
	select {
	case <-timer.C:
		log.Println("Timeout! Not Found!")
		cancel()
	case <-resCh:
		log.Println("Found it!")
		cancel()
	}
}

func SearchTarget(ctx context.Context, data []int, target int, ch chan int) {
	for _, v := range data {
		select {
		case <-ctx.Done():
			log.Println("Task cancelded!")
			return
		default:
		}
		time.Sleep(time.Millisecond * 1500)
		if target == v {
			ch <- v
			return
		}
	}
}
