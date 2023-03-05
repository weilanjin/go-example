package goroutine_test

import (
	"fmt"
	"runtime"
	"sync"
	"testing"

	"lovec.wlj/pkg/mem"
)

// 闭包
func TestClosure(t *testing.T) {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func(salutation string) {
			wg.Done()
			fmt.Println(salutation)
		}(salutation)
	}
	wg.Wait() // 只是一个 join 点
}

func TestGoroutineSize(t *testing.T) {

	// 利用 goroutine 泄露来测试 goroutine 大小
	var c <-chan any
	var wg sync.WaitGroup
	noop := func() { wg.Done(); <-c }

	const numGoroutines = 1e4
	wg.Add(numGoroutines)

	before := mem.MemConsumed()
	for i := numGoroutines; i > 0; i-- {
		go noop()
	}
	wg.Wait()
	after := mem.MemConsumed()
	fmt.Println(runtime.NumGoroutine())
	fmt.Printf("%.3fkb\n", float64(after-before)/numGoroutines/1000)
}

// go test -bench=.-cpu=1 -v ./gorotine_test.go
func BenchmarkContextSwitch(b *testing.B) {
	var wg sync.WaitGroup
	begin := make(chan struct{})
	c := make(chan struct{})

	var token struct{}
	sender := func() {
		defer wg.Done()
		<-begin
		for i := 0; i < b.N; i++ {
			c <- token
		}
	}
	receiver := func() {
		defer wg.Done()
		<-begin
		for i := 0; i < b.N; i++ {
			<-c
		}
	}
	wg.Add(2)
	go sender()
	go receiver()
	b.StartTimer()
	close(begin) // 关闭阻塞
	wg.Wait()
}
