package limit

import (
	"sync"
	"time"
)

// 令牌桶
// 允许流量突发
type TokenBucket struct {
	rate         int64
	capacity     int64 // 桶的容量
	tokens       int64 // 桶中当前token数据
	lastTokenSec int64 // 桶上次放token的时间戳 s

	mu sync.Mutex
}

func (l *TokenBucket) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now().Unix()
	l.tokens = l.tokens + (now-l.lastTokenSec)*l.rate // 先添加令牌
	if l.tokens > l.capacity {
		l.tokens = l.capacity
	}
	l.lastTokenSec = now
	if l.tokens > 0 {
		// 还有令牌，领取令牌
		l.tokens--
		return true
	}
	// 没有令牌，则拒绝
	return false
}

func (l *TokenBucket) Set(r, c int64) {
	l.rate = r
	l.capacity = c
	l.tokens = 0
	l.lastTokenSec = time.Now().Unix()
}
