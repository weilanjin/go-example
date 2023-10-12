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

var mu sync.Mutex
var chain string

// fatal error: all goroutines are asleep - deadlock!
// 使用Lock()加锁之后, 不能再继续对其加锁. 直到利用Unlock()解锁后才能再加锁
func A2() {
	mu.Lock()
	defer mu.Unlock()
	chain = chain + " --> A"
	B2()
}

func B2() {
	chain = chain + " -->B"
	C2()
}

func C2() {
	mu.Lock()
	defer mu.Unlock()
	chain = chain + " -->C"
}

func TestMutex1(t *testing.T) {
	chain = "main"
	A2()
	fmt.Println(chain)
}

type MyMutex struct {
	count int
	sync.Mutex
}

func TestMutex2(t *testing.T) {
	var mu MyMutex
	mu.Lock()
	var mu1 = mu // mu1 现在是加状态
	mu.count++
	mu.Unlock()
	// 重复加锁
	mu1.Lock() // fatal error: all goroutines are asleep - deadlock!
	mu1.count++
	mu1.Unlock()
	fmt.Println(mu.count, mu1.count)
}
