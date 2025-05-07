package question

import (
	"fmt"
	"testing"
)

// go 多重赋值，从右到左顺序
func TestMultipleAssignments(t *testing.T) {
	i := 1
	s := []string{"a", "b", "c"}
	i, s[i-1] = 2, "z"
	fmt.Println(s) // [z b c]
}