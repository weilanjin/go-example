package mutex

import "sync"

// debug
func copyMutex() {
	var mu sync.Mutex  // 第一个锁
	var mu2 sync.Mutex // 第二个锁

	mu.Lock() // 第一个加锁
	defer mu.Unlock()

	mu2 = mu   // 把第一个锁复制给第二个锁, 第二个锁处于锁的状态
	mu2.Lock() // 阻塞在这里
	// ----
	mu2.Unlock()

}
