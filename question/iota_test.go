package question

import (
	"fmt"
	"testing"
)

// https://studygolang.com/articles/2192
// 1. iota 在 const 关键字出现时将被重置为0。
// 2.const中每新增一行常量声明将使 iota 计数一次。

const (
	a = iota
	b = iota
)

const (
	x = iota
	_
	y
	z = "zz"
	k
	p = iota
)

func TestIota(t *testing.T) {
	fmt.Println(a, b) // 0, 1
	fmt.Println(x, y, z, k, p)
}

// output:
// 0, 2 zz, zz, 5