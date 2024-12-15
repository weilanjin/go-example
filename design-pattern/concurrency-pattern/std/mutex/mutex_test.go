package mutex

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// Mutex
// Lock() 获得锁
// UnLock() 释放锁 [1.如果直接释放一个未加锁的Mutex, 则会panic] sync: unlock of unlocked mutex
// TryLock() 尝试获得锁 Go 1.18 +

func TestUnLock(t *testing.T) {
	var tags struct {
		mu     sync.Mutex
		number int
	}

	go func() {
		fmt.Println("go func 1")
		tags.mu.Unlock() // fatal("sync: unlock of unlocked mutex")
		time.Sleep(time.Second)
		tags.number = 1
	}()
	time.Sleep(1500 * time.Millisecond)
	fmt.Printf("number: %v\n", tags.number)
}
