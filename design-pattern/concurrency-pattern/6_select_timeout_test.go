package concurrency_pattern

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

func boringChannel(msg string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond)
		}
	}()
	return c
}

func Test6(t *testing.T) {
	c := boringChannel("Job")
	timeout := time.After(5 * time.Second)
	for {
		select {
		case s := <-c:
			log.Println(s)
		case <-timeout:
			log.Println("You talk to much.")
			return
		}
	}
}
