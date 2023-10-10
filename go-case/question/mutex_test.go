package question

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type data struct {
	sync.Mutex
}

func (d data) test(s string) {
	// 锁失效
	// 将 Mutex 作为匿名字段时，相关的方法必须使用指针接收者，否则会导致锁机制失效
	d.Lock()
	defer d.Unlock()
	for i := 0; i < 5; i++ {
		fmt.Println(s, i)
		time.Sleep(time.Second)
	}
}

func TestMutex(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	var d data
	go func() {
		defer wg.Done()
		d.test("read")
	}()
	go func() {
		defer wg.Done()
		d.test("write")
	}()
	wg.Wait()
}