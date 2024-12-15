package rwmutex

import (
	sync "github.com/sasha-s/go-deadlock" //
	// "sync"
	"testing"
	"time"
)

// 测试死锁
// github.com/sasha-s/go-deadlock 可以分析死锁的位置

func TestDeadLock(t *testing.T) {
	// 分析 reader 先调用读锁，然后调用写锁，导致死锁

	var mu sync.RWMutex
	go func() {
		mu.RLock() // 1.获取读锁
		{
			time.Sleep(5 * time.Second)
			mu.RLock() // 2.再次获取读锁
			t.Log("程序不可能执行到这里")
			mu.RUnlock()
		}
		mu.RUnlock()
	}()
	time.Sleep(1 * time.Second)
	mu.Lock() // 3. 2 前获取写锁
	t.Log("程序不可能执行到这里")
	mu.Unlock()
}
