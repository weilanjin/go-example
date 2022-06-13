package concurrency_pattern

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

func boringSelect(msg string, quit chan string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			select {
			case c <- fmt.Sprintf("%s %d", msg, i):
			case <-quit:
				log.Println("clean up")
				quit <- "See you!"
				return
			}
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return c
}

func Test7(t *testing.T) {
	quit := make(chan string)
	c := boringSelect("Joe", quit)
	for i := 3; i >= 0; i-- {
		log.Println(<-c)
	}
	quit <- "Bye"
	log.Println("Joe sey:", <-quit)
}
