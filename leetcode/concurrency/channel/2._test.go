package channel

import (
	"log"
	"testing"
)

func Test2(t *testing.T) {
	var ch chan int
	var count int
	go func() {
		ch <- 1
	}()
	go func() {
		count++
		close(ch) // close of nil channel
	}()
	<-ch
	log.Println(count)
}
