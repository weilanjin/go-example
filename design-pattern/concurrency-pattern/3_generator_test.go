package concurrency_pattern

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

func generator(msg string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; i < 10; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
		close(c)
	}()
	return c
}

func Test3(t *testing.T) {
	joe := generator("Joe")
	ahn := generator("Ahn")
	for i := 0; i < 10; i++ {
		log.Println(<-joe)
		log.Println(<-ahn)
	}
	log.Println("You're both generator, I'm leaving")
}
