package question

import (
	"fmt"
	"testing"
)

func TestSelect(t *testing.T) {
	c := make(chan int, 1)
	for range [3]struct{}{} {
		select {
		default:
			fmt.Println(1)
		case <-c:
			fmt.Println(2)
			c = nil
		case c <- 1:
			fmt.Println(3)
		}
	}
}

// output:
// 3 2 1