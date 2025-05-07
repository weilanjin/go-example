package question

import (
	"fmt"
	"testing"
)

func TestVariableParam(t *testing.T) {
	fn := func(num ...int) {
		num[0] = 18
	}

	arr := []int{5, 6, 7}
	fn(arr...)
	fmt.Println(arr[0]) // 18
}