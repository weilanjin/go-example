package concurrency_test

import (
	"fmt"
	"testing"
	"time"
)

// for { // Either loop infinitely or range over something
// 		select {
// 		// Do some work with channels
// 		}
//  }

func TestForSelect(t *testing.T) {
	done := make(chan any)
	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	workCounter := 0
loop:
	for {
		select {
		case <-done:
			break loop
		default:
		}
		workCounter++
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("Achieved %v cycles of work before signalled to stop.\n", workCounter)
}

type foo int
type bar int

func TestXxx(t *testing.T) {
	m := make(map[any]int)
	m[foo(1)] = 2
	m[bar(1)] = 3
	fmt.Println(m)
}
