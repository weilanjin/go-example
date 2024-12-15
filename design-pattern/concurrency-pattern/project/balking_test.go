package project

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 回避(balking)
// 当一个对象尝试执行某个操作时, 如果发现当前的状态不适合执行该操作,它会停止执行,而不是继续执行下去
// - 在多线程编程中,如果某个线程发现共享变量的状态已经发生改变,那么它可能需要停止执行某个操作,以避免出现竞态条件.
// - 在编写网络应用程序时,如果客户端向服务器发送请求时发现网络连接已经中断,那么它可以立即停止发送请求,而不是等待网络连接重新建立后再尝试发送请求.

func TestBalking(t *testing.T) {
	var flag atomic.Bool

	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				if !flag.CompareAndSwap(false, true) { // 已经有 goroutine 在执行了,回避
					time.Sleep(time.Second)
					t.Logf("balking %d - %d", id, i)
					continue
				}

				// 执行业务逻辑

				flag.Store(false)
			}
		}(i)
	}
	wg.Wait()
}