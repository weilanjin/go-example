package mutex

import (
	"log"
	"sync"
	"testing"
)

var mu sync.Mutex
var chain string

func Test1(t *testing.T) {
	chain = "main"
	A()
	log.Println(chain)
}

func A() {
	mu.Lock()
	defer mu.Unlock()
	chain = chain + " --> A"
	B()
}

func B() {
	chain = chain + " --> B"
	C()
}

func C() {
	mu.Lock() // fatal error: all goroutines are asleep - deadlock!
	defer mu.Unlock()
	chain = chain + " --> C"
}
