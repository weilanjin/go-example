package concurrency_pattern

import (
	"log"
	"testing"
)

func f(left, right chan int) {
	left <- 1 + <-right
}

func Test8(t *testing.T) {
	const n = 1000
	leftMost := make(chan int)
	left := leftMost
	right := leftMost

	for i := 0; i < n; i++ {
		right = make(chan int)
		go f(left, right)
		left = right
	}
	go func(c chan int) {
		c <- 1
	}(right)
	log.Println(<-leftMost)
}
