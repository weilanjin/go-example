package question

import (
	"fmt"
	"testing"
	"time"
)

func TestChan(t *testing.T) {
	ch := make(chan int, 100)
	go func() {
		for i := 1; i < 10; i++ {
			ch <- i
		}
	}()
	go func() {
		for {
			a, ok := <-ch
			if !ok {
				fmt.Println("close")
				return
			}
			fmt.Println(a)
		}
	}()
	// 第一个 goroutine 还没有写完就被 close
	close(ch) // panic: send on closed channel
	fmt.Println("ok")
	time.Sleep(time.Second * 10)
}