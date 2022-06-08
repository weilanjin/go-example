package mutex

import (
	"log"
	"sync"
	"testing"
)

type MyMutex struct {
	count int
	sync.Mutex
}

func Test3(t *testing.T) {
	var mu MyMutex
	mu.Lock()
	var mu2 = mu
	mu.count++
	mu.Unlock()

	mu2.Lock() // fatal error: all goroutines are asleep - deadlock!
	mu2.count++
	mu2.Unlock()
	log.Println(mu.count, mu2.count)
}
