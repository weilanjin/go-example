package question

import (
	"fmt"
	"testing"
)

func foo(x interface{}) {
	if x == nil {
		fmt.Println("empty interface")
		return
	}
	fmt.Println("non-empty interface")
}

func TestInterfaceNil(t *testing.T) {
	foo(nil) // empty interface
	// interface 的内部结构，接口除了有静态类型，还有动态类型和动态值，当且仅当动态值和动态类型都为nil时，接口类型值才为nil。
	var x *int = nil
	foo(x) // non-empty interface
}

func TestInterfaceAssert(t *testing.T) {
	x := interface{}(nil)
	y := (*int)(nil)
	a := y == x
	b := y == nil
	// 类型断言：
	// i.(Type), i 是接口，Type是类型或接口
	_, c := x.(interface{})
	fmt.Println(a, b, c) // false true false
}

func TestInterfaceCompare(t *testing.T) {
	var x interface{}
	var y interface{} = []int{3, 5}
	_ = x == x
	_ = x == y
	_ = y == y // runtime error: comparing uncomparable type []int [recovered]
}