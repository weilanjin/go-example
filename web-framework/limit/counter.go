package limit

import (
	"sync"
	"time"
)

// 引用计数器使用到了锁
// 在高并发场景不太实用
type Counter struct {
	rate  int           // 计数周期内最多允许的请求数
	begin time.Time     // 计数开始时间
	cycle time.Duration // 计数周期
	count int           // 计数周期内累计收到的请求数
	mu    sync.Mutex
}

//
func (l *Counter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.count == l.rate-1 {
		now := time.Now()
		if now.Sub(l.begin) >= l.cycle {
			l.Reset(now) // 速度在允许范围内，重置计数器
			return true
		}
		return false
	}
	// 没有达到速率限制，计数 +1
	l.count++
	return true
}

func (l *Counter) Reset(now time.Time) {
	l.begin = now
	l.count = 0
}

func (l *Counter) Set(rate int, cycle time.Duration) {
	l.rate = rate
	l.begin = time.Now()
	l.cycle = cycle
	l.count = 0
}
