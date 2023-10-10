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