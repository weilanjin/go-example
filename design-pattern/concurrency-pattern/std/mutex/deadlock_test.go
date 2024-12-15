package mutex

import (
	"net/http"
	_ "net/http/pprof"
	"sync"
	"testing"
)

// 浏览器访问 debug=2 可以详细的看到每一个goroutine的堆站信息和状态
// http://localhost:8080/debug/pprof/goroutine?debug=2
// goroutine 都处于 semacquire 状态
// semacquire 这些携程在等待一个信号或资源释放后才能继续执行.

/*
goroutine 36 [sync.Mutex.Lock]:
sync.runtime_SemacquireMutex(0x0?, 0x0?, 0x0?)
	/usr/local/go/src/runtime/sema.go:77 +0x28
sync.(*Mutex).lockSlow(0x140001a22e8)
	/usr/local/go/src/sync/mutex.go:171 +0x330
sync.(*Mutex).Lock(0x140001a22e8)
	/usr/local/go/src/sync/mutex.go:90 +0x94
lovec.wlj/design-patten/concurrency-pattern/concurrency-program/mutex.TestDeadlock.func1()
	/Users/lanjin/Documents/work/code/go-example/design-pattern/concurrency-pattern/concurrency-program/mutex/deadlock_test.go:15 +0x30
created by lovec.wlj/design-patten/concurrency-pattern/concurrency-program/mutex.TestDeadlock in goroutine 35
	/Users/lanjin/Documents/work/code/go-example/design-pattern/concurrency-pattern/concurrency-program/mutex/deadlock_test.go:14 +0xf8

goroutine 37 [sync.Mutex.Lock]:
sync.runtime_SemacquireMutex(0x0?, 0x0?, 0x0?)
	/usr/local/go/src/runtime/sema.go:77 +0x28
sync.(*Mutex).lockSlow(0x140001a22e8)
	/usr/local/go/src/sync/mutex.go:171 +0x330
sync.(*Mutex).Lock(0x140001a22e8)
	/usr/local/go/src/sync/mutex.go:90 +0x94
lovec.wlj/design-patten/concurrency-pattern/concurrency-program/mutex.TestDeadlock.func1()
	/Users/lanjin/Documents/work/code/go-example/design-pattern/concurrency-pattern/concurrency-program/mutex/deadlock_test.go:15 +0x30
created by lovec.wlj/design-patten/concurrency-pattern/concurrency-program/mutex.TestDeadlock in goroutine 35
	/Users/lanjin/Documents/work/code/go-example/design-pattern/concurrency-pattern/concurrency-program/mutex/deadlock_test.go:14 +0xf8
*/

func TestDeadlock(t *testing.T) {
	var count int64
	var mu sync.Mutex
	for i := 0; i < 100; i++ {
		go func() {
			mu.Lock()
			// defer mu.Unlock() // 认为的故意不释放
			count++
		}()
	}
	if err := http.ListenAndServe(":8080", nil); err != nil {
		t.Fatal(err)
	}
}
