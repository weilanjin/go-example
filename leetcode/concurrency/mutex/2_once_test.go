package mutex

import (
	"log"
	"sync"
	"testing"
	"time"
)

func Test2(t *testing.T) {
	o := &Once{}
	for i := 0; i < 10; i++ {
		go o.Do(func() {
			log.Println(123)
		})
	}
	time.Sleep(time.Second * 4)
}

// 双检查实现单例
type Once struct {
	mu   sync.Mutex
	done uint32
}

func (o *Once) Do(f func()) {
	if o.done == 1 {
		return
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.done == 0 {
		o.done = 1
		f()
	}
}
