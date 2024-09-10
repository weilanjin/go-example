package rwmutex

import (
	"sync"
	"testing"
)

// go test -run ^$ -bench BenchmarkCounter

// 读写的大概比例是 10000:1
// goos: darwin
// goarch: arm64
// pkg: lovec.wlj/design-patten/concurrency-pattern/concurrency-program/rwmutex
// BenchmarkCounterMutex-10    	   10000	   1224068 ns/op
func BenchmarkCounterMutex(b *testing.B) {
	var counter int64
	var mu sync.Mutex
	for i := 0; i < b.N; i++ {
		b.RunParallel(func(p *testing.PB) {
			i := 0
			for p.Next() {
				i++
				if i%10000 == 0 { // 10000的整数倍时一次加锁
					mu.Lock()
					counter++
					mu.Unlock()
				} else { // 只读
					mu.Lock()
					_ = counter
					mu.Unlock()
				}
			}
		})
	}
}

// - 需要对写进行保护时, 调用写锁
// - 需要对读进行保护时, 调用读锁
// goos: darwin
// goarch: arm64
// pkg: lovec.wlj/design-patten/concurrency-pattern/concurrency-program/rwmutex
// BenchmarkCounterRWMutex-10    	   10000	   1206542 ns/op
func BenchmarkCounterRWMutex(b *testing.B) {
	var counter int64
	var mu sync.RWMutex
	for i := 0; i < b.N; i++ {
		b.RunParallel(func(p *testing.PB) {
			i := 0
			for p.Next() {
				i++
				if i%10000 == 0 { // 10000的整数倍时一次加锁
					mu.Lock()
					counter++
					mu.Unlock()
				} else { // 只读
					mu.RLock()
					_ = counter
					mu.RUnlock()
				}
			}
		})
	}
}
