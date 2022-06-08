package waitgroup

import (
	"sync"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		time.Sleep(time.Millisecond)
		wg.Done()
		wg.Add(1)
	}()
	// wait() 和 add() 不能并发
	wg.Wait() // panic: sync: WaitGroup is reused before previous Wait has returned
}
