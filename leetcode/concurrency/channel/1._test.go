package channel

import (
	"log"
	"runtime"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	var ch chan int
	go func() {
		ch = make(chan int, 1)
		ch <- 1
	}()
	go func(ch chan int) {
		time.Sleep(time.Second)
		<-ch
	}(ch)
	c := time.Tick(1 * time.Second)
	for range c {
		log.Printf("#goroutines: %d\n", runtime.NumGoroutine()) // 3
	}
}
