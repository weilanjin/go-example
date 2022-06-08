package rwmutex

import (
	"log"
	"sync"
	"testing"
	"time"
)

var mu sync.RWMutex
var count int

func Test1(t *testing.T) {
	go A()
	time.Sleep(2 * time.Second)
	mu.Lock() // fatal error: all goroutines are asleep - deadlock!
	defer mu.Unlock()
	count++
	log.Println(count)
}

func A() {
	mu.RLock()
	defer mu.RUnlock()
	B()
}

func B() {
	time.Sleep(5 * time.Second)
	C()
}

func C() {
	mu.RLock()
	defer mu.RUnlock()
}
