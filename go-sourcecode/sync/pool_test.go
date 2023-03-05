package sync_test

import (
	"fmt"
	"sync"
	"testing"
)

func TestPool(t *testing.T) {
	pool := &sync.Pool{
		New: func() any {
			fmt.Println("creating new instance.")
			return struct{}{}
		},
	}

	pool.Get()
	o := pool.Get()
	pool.Put(o)
	pool.Get()
	// output:
	// creating new instance.
	// creating new instance.
}

func TestPoolCalc(t *testing.T) {
	var numCalcsCreated int
	calcPool := &sync.Pool{
		New: func() any {
			numCalcsCreated += 1
			mem := make([]byte, 1024)
			return &mem
		},
	}

	// Seed the pool with 4KB
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())

	const numWorkers = 1024 * 10224
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := numWorkers; i > 0; i-- {
		go func() {
			defer wg.Done()
			mem := calcPool.Get().(*[]byte)
			defer calcPool.Put(mem)
		}()
	}

	wg.Wait()
	fmt.Printf("%d calculators were created.", numCalcsCreated)
}
