package question

import (
	"fmt"
	"testing"
)

type N int

func (n N) test() {
	fmt.Println(n)
}

func TestMethod(t *testing.T) {
	var n N = 10
	fmt.Println(n) // 10

	// 方法表达是。
	// 通过类型引用的方法表达式会被还原成普通函数样式，接收者是第一个参数，调用时显示传参。
	n++
	N.test(n) // 11

	n++
	(*N).test(&n) // 12
}

func (n *N) test1() {
	fmt.Println(*n)
}

func TestMethod1(t *testing.T) {
	var n N = 10
	p := &n

	n++
	f1 := n.test1

	n++
	f2 := p.test1

	n++
	fmt.Println(n) // 13
	// 当目标方法接收者是指针类型时，那么复制的就是指针
	f1() // 13
	f2() // 13
}