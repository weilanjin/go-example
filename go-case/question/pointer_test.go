package question

import (
	"fmt"
	"testing"
)

type P *int
type Q *int

func TestPointer(t *testing.T) {
	var p P = new(int)
	*p += 8
	var x *int = p
	var q Q = x
	*q++
	// 变量指针指向相同的地址
	fmt.Println(*p, *q) // 9 9
}
