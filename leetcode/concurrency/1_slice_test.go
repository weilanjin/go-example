package concurrency

import (
	"log"
	"sync"
	"testing"
)

func Test1(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	var ints = make([]int, 0, 1000)
	go func() {
		for i := 0; i < 1000; i++ {
			ints = append(ints, i)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 1000; i++ {
			ints = append(ints, i)
		}
		wg.Done()
	}()
	wg.Wait()
	log.Println(len(ints))
}
