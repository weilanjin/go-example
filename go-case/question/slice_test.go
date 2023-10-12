package question

import (
	"fmt"
	"testing"
)

// 数组或切片的截取操作
// 带2个或者3个参数 [i:j] 和 [i:j:k]
// 1. i 省略，默认0
// 2. j 省略。 默认底层数组长度
// 3. k 主要是用于来限制切片的容量，但是不能大于数组的长度
func TestSlice(t *testing.T) {
	s := [3]int{1, 2, 3}
	a := s[:0]         // 0, 3
	b := s[:2]         // 2, 3
	c := s[1:2:cap(s)] // 1, 2
	fmt.Println(len(a), len(b), len(c), cap(a), cap(b), cap(c))
}

func TestArr(t *testing.T) {
	s1 := []int{1, 2, 3}
	s2 := s1[1:]
	s2[1] = 4
	// s2 底层引用 s1，所以修改值会改变也会改变 s1
	fmt.Println(s1) // [1 2 4]
	s2 = append(s2, 5, 6, 7)
	// append 会发生扩容，s2 引用了新的数组
	fmt.Println(s1) // [1 2 4]
}

func TestSliceIndex(t *testing.T) {
	// 字面量初始切片时候，可以指定索引，没有指定索引的元素会在前一个索引基础之上+1。
	var x = []int{2: 2, 3, 0: 1}
	fmt.Println(x)
}

func TestSlice3(t *testing.T) {
	s := make([]int, 3, 9)
	fmt.Println(len(s)) // 3
	_ = s[6:]           // 报错： [6:3] [i:j] 如果 j 省略 默认len的长度

	// 从一个基础切片派生出的子切片的长度可能大于基础切片的长度
	// [low, high] 0 <= low <= high <= cap
	s2 := s[4:8]
	fmt.Println(len(s2)) // 4
}

func TestSlice4(t *testing.T) {
	a := [3]int{0, 1, 2}
	s := a[1:2]

	s[0] = 11
	s = append(s, 12)
	s = append(s, 13)
	s[0] = 21

	fmt.Println(a) // [0, 11, 12]
	fmt.Println(s) // [21, 12, 13]
}

func TestSlice5(t *testing.T) {
	var src, dst []int
	src = []int{1, 2, 3}
	// dst = make([]int, len(src))
	// 或者
	// dst = append(src, src...)
	copy(dst, src)
	// dst 必须分配足够的空间
	fmt.Println(dst) // []
}

func TestSlice6(t *testing.T) {
	s := []int{0, 1}
	// 对一个切片执行[i,j]的时候, i 和 j都不能超过切片的长度值
	fmt.Println(s[len(s):]) // []
}
