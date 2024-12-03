package memmodel

import "sync"

// 第n次的UnLock调用一定synchronized before 第n次Lock方法的返回

func Mutex() {
	var mu sync.Mutex
	var s string

	var foo = func() {
		s = "hello world"
		mu.Unlock()
	}
	mu.Lock()
	go foo()
	mu.Lock()
	print(s) // hello world
}