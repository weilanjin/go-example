package question

import (
	"fmt"
	"testing"
	"time"
)

// 创建每个元素的副本，而不是元素的引用
func TestFor(t *testing.T) {
	slice := []int{0, 1, 2, 3}
	m := make(map[int]*int)

	for i, v := range slice {
		// 解决
		// value := v
		// m[i] = value
		m[i] = &v
	}

	// 输出
	for k, v := range m {
		fmt.Println(k, "-->", *v)
	}
}

// output :
// 		2 --> 3
// 		3 --> 3
//  	0 --> 3
// 		1 --> 3

func TestForSlice(t *testing.T) {
	v := []int{1, 2, 3}
	// 循环次数在循环开始前就已经确定，循环内改变切⽚的⻓度，不影响循环次数
	for i := range v {
		v = append(v, i)
	}
	fmt.Println(v) // [1 2 3 0 1 2]
}

func TestSlice1(t *testing.T) {
	var m = [...]int{1, 2, 3}
	// 使⽤短变量声明 (:=) 的形式迭代变量，需要注意的是，变量 i、v 在每次循环体中都会被重⽤，⽽不是重新声明
	for i, v := range m {
		go func() {
			fmt.Println(i, v)
		}()
	}
	time.Sleep(3 * time.Second)
}

// output:
//		2 3
//		2 3
//		2 3

func TestSlice2(t *testing.T) {
	var a = [5]int{1, 2, 3, 4, 5} // fix: 第一种：a 换成slice
	var r [5]int
	// range a 是 var a 的副本， 修改 var a 不会改变 range
	// fix： 第二种 for i, v := range &a {
	for i, v := range a {
		if i == 0 {
			a[1] = 12
			a[2] = 13
		}
		r[i] = v
	}
	fmt.Println(a) // [1 12 13 4 5]
	fmt.Println(r) // [1 2 3 4 5]
}

func TestFor1(t *testing.T) {
	var k = 9
	for k = range []int{} {
	}
	fmt.Println(k) // 9

	for k = 0; k < 3; k++ {
	}
	fmt.Println(k)                 // 3
	for k = range (*[3]int)(nil) { // k is index
	}
	fmt.Println(k) // 2
}