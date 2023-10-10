package question

import (
	"fmt"
	"testing"
)

func TestAppend(t *testing.T) {
	s := make([]int, 5)
	s = append(s, 1, 2, 3)
	// len = 5, 前面补零，往后追加
	fmt.Println(s) // [0 0 0 0 0 1 2 3]

	s1 := make([]int, 0)
	s1 = append(s1, 1, 2, 3)
	fmt.Println(s1) // [1 2 3]
}

// -------------------
func change(s ...int) {
	s = append(s, 3)
}

func TestVarargs(t *testing.T) {
	slice := make([]int, 5, 5) // [0, 0, 0, 0, 0]
	slice[0] = 1
	slice[1] = 2          // [1, 2, 0, 0, 0]
	change(slice...)      // 产生了扩容 s 没有引用 slice
	fmt.Println(slice)    // [1, 2, 0, 0, 0]
	change(slice[0:2]...) // [1, 2] 容量是 5， 没有产生扩容
	fmt.Println(slice)    // [1, 2, 3, 0, 0] // 输出原始的slice
}