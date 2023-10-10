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