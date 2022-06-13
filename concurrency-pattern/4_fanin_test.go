package concurrency_pattern

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

func boringFanin(msg string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return c
}

func fanin(cs ...<-chan string) <-chan string {
	c := make(chan string)
	for _, ci := range cs {
		go func(cv <-chan string) {
			for {
				c <- <-cv
			}
		}(ci)
	}
	return c
}

func Test4(t *testing.T) {
	c := fanin(boringFanin("Joe"), boringFanin("Ahn"))
	for i := 0; i < 5; i++ {
		log.Println(<-c)
	}
	log.Println("You're both boring. I'm leaving")
}
