package stream

import (
	"fmt"
	"math"
	"testing"
)

func TestStream(t *testing.T) {
	Of(1, 2, 3).
		Filter(func(i any) bool {
			return i.(int) > 1
		}).
		Map(func(i any) any {
			return "hello" + fmt.Sprintf("%d", i)
		}).
		ForEach(func(s any) {
			fmt.Println(s)
		})

	fmt.Println("+", Of(1, 2, 3).Reduce(func(x any, y any) any { return x.(int) + y.(int) }, 4))
	fmt.Println("min", Of(1, 2, 3).Reduce(func(x any, y any) any { return min(x.(int), y.(int)) }, math.MaxInt))
	fmt.Println("max", Of(1, 2, 3).Reduce(func(x any, y any) any { return max(x.(int), y.(int)) }, math.MinInt))

	res := Of(1, 2, 2, 3, 4, 5, 5, 5, 6, 9).Collect(func() any {
		return make([]int, 0)
	}, func(s any, i any) {
		s = append(s.([]int), i.(int))
	})

	fmt.Println(res)
}
