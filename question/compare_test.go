package question

import (
	"fmt"
	"testing"
)

type foo1 struct{ Val int }
type bar struct{ Val int }

func TestCompare(t *testing.T) {
	a := &foo1{Val: 5}
	b := &foo1{Val: 5}
	c := foo1{Val: 5}
	d := bar{Val: 5}
	e := bar{Val: 5}
	// go 没有引用变量,每个变量都占有一个唯一的内存位置
	// false true true
	fmt.Println(a == b, c == foo1(d), d == e)
}
