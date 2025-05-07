package question

import (
	"fmt"
	"testing"
)

// recover() 必须在defer()函数中调用才有效

func TestRecover(t *testing.T) {
	defer func() {
		fmt.Println(recover()) // 捕获 panic(1)
	}()
	defer func() {
		defer fmt.Println(recover()) // 捕获 panic(2)
		panic(1)
	}()
	defer recover() // 无效捕获
	panic(2)
}

// output:
//       2 1

func fp() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	panic(1) // 下面的 panic 不会执行
	panic(2)
}
