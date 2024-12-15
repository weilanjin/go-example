package goroutine

import (
	"fmt"
	"testing"
	"time"
)

// 函数是异步执行的,函数异步执行时参数会在后面某个时间才被使用,
// 传入的参数调用者所在的 goroutine 中可能会被修改.

func TestGoFunc(t *testing.T) {
	list := []int{1}

	foo := func(l int) {
		time.Sleep(1 * time.Second)
		fmt.Printf("passed len: %d, current list len %d\n", l, len(list))
	}

	go foo(len(list)) // 获取到cpu资源就执行 list(slice) 是一个引用类型

	list = append(list, 2)

	foo = func(l int) {
		fmt.Printf("passed len: %d, current list len %d\n", l*100, len(list)*100)
	}
	time.Sleep(2 * time.Second)
	foo(len(list))

	// output
	// passed len: 1, current list len 2
	// passed len: 200, current list len 200
}
