package question

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

func TestHunger(t *testing.T) {
	runtime.GOMAXPROCS(1)
	go func() {
		for i := 0; i < 10; i++ {
			println(i)
		}
	}()
	// for {} 会独占cpu资源导致其他 goroutine 饿死
	// select {}
	for {
	}
}

func TestGosched(t *testing.T) {
	runtime.GOMAXPROCS(1)

	var N = 26
	var wg sync.WaitGroup
	wg.Add(2 * N)
	for i := 0; i < N; i++ {
		go func(i int) {
			defer wg.Done()
			runtime.Gosched() // 让出协程调度
			fmt.Printf("%c", 'a'+i)
		}(i)
		go func(i int) {
			defer wg.Done()
			fmt.Printf("%c", 'A'+i)
		}(i)
	}
	wg.Wait()
}
