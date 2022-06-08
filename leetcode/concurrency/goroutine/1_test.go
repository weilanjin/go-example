package goroutine

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type T struct {
	V int
}

func (t *T) Incr(wg *sync.WaitGroup) {
	t.V++
	wg.Done()
}

func (t *T) Print() {
	time.Sleep(1e9)
	fmt.Print(t.V) // 9999999999
	fmt.Println()
}

func Test1(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(10)
	var ts = make([]T, 10)
	for i := 0; i < 10; i++ {
		ts[i] = T{i}
	}
	for _, t := range ts {
		go t.Incr(&wg)
	}
	wg.Wait()

	for _, t := range ts {
		go t.Print()
	}
	time.Sleep(6 * time.Second)
}
