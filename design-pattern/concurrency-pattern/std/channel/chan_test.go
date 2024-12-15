package channel

import (
	"fmt"
	"testing"
	"time"
)

func TestSendNilChan(t *testing.T) {
	var ch chan int
	ch <- 1 // fatal error: all goroutines are asleep - deadlock!
	for v := range ch {
		t.Log(v)
	}
}

func TestNewMutex(t *testing.T) {
	// 创建一个新的互斥锁
	mutex := NewMutex()

	// 模拟多个goroutine竞争锁
	go func() {
		if mutex.TryLock() {
			fmt.Println("Goroutine 1 acquired the lock")
			time.Sleep(2 * time.Second)
			mutex.UnLock()
			fmt.Println("Goroutine 1 released the lock")
		} else {
			fmt.Println("Goroutine 1 failed to acquire the lock")
		}
	}()

	go func() {
		if mutex.LockTimeout(1 * time.Second) {
			fmt.Println("Goroutine 2 acquired the lock")
			time.Sleep(1 * time.Second)
			mutex.UnLock()
			fmt.Println("Goroutine 2 released the lock")
		} else {
			fmt.Println("Goroutine 2 failed to acquire the lock within timeout")
		}
	}()

	// 主goroutine尝试获取锁
	mutex.Lock()
	fmt.Println("Main goroutine acquired the lock")
	time.Sleep(3 * time.Second)
	mutex.UnLock()
	fmt.Println("Main goroutine released the lock")

	// 等待所有goroutine完成
	time.Sleep(5 * time.Second)

	// output:
	// Main goroutine acquired the lock
	// Goroutine 1 failed to acquire the lock
	// Goroutine 2 failed to acquire the lock within timeout
	// Main goroutine released the lock
}

func TestTaskOr(t *testing.T) {
	sig := func(after time.Duration) <-chan any {
		c := make(chan any)
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-OrV2(
		sig(10*time.Second),
		sig(20*time.Second),
		sig(30*time.Second),
		sig(2*time.Second),
		sig(40*time.Second),
		sig(50*time.Second),
	)
	fmt.Printf("done after %v\n", time.Since(start))

	// output:
	// done after 2.001299542s
}