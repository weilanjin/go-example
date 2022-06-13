package concurrency_pattern

import (
	"log"
	"testing"
	"time"
)

type Ball struct {
	hits int
}

func player(name string, table chan *Ball) {
	for {
		ball := <-table
		ball.hits++
		log.Println(name, ball.hits)
		time.Sleep(100 * time.Millisecond)
		table <- ball
	}
}

func Test10(t *testing.T) {
	table := make(chan *Ball)
	go player("ping", table)
	go player("pong", table)
	table <- new(Ball)
	time.Sleep(1 * time.Second)
	<-table
	panic("show me the stack")
}
