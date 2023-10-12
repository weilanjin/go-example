package question

import (
	"fmt"
	"testing"
	"time"
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

// ------------------------------
func A() int {
	time.Sleep(100 * time.Millisecond)
	return 1
}

func B() int {
	time.Sleep(1000 * time.Millisecond)
	return 2
}

func TestSelect1(t *testing.T) {
	ch := make(chan int, 1)
	go func() {
		select {
		case ch <- A():
		case ch <- B():
		default:
			ch <- 3
		}
	}()
	fmt.Println(<-ch)
}

// output:
//   1 或 2 随机出
