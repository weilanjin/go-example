package concurrency_pattern

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

func boringCh(msg string, c chan string) {
	for i := 0; ; i++ { // go 会泄漏
		// send the value to channel
		// it also waits for receiver to be ready
		c <- fmt.Sprintf("%s %d", msg, i)
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}
}

func Test2(t *testing.T) {
	c := make(chan string)
	go boringCh("boring", c)
	for i := 0; i < 5; i++ {
		log.Printf("You say: %q\n", <-c)
	}
	log.Println("You're boring, I'm leaving")
}
