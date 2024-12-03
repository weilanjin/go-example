package memmodel

import "sync"

// 对于 once.Do(f) 调用, f函数的单次调用一定 synchronized before 任何 once.Do(f) 调用的返回

func onceDo() {
	var s string
	var once sync.Once

	var foo = func() {
		s = "hello world"
	}

	once.Do(foo)
	print(s)
}