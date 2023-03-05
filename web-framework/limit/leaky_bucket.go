package limit

import (
	"math"
	"sync"
	"time"
)

// 漏桶算法
// 水溢出请求拒绝
// 不适应流量突发的情况
type LeakyBucket struct {
	rate       float64 // 固定每秒出水速率
	capacity   float64 // 桶的容量
	water      float64 // 桶中当前数量
	lastLeakMs int64   // 桶上次漏水时间戳 ms
	mu         sync.Mutex
}

func (l *LeakyBucket) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now().UnixNano() / 1e6
	eclipse := float64(now-l.lastLeakMs) * l.rate / 1000
	l.water = l.water - eclipse
	l.water = math.Max(0, l.water)
	l.lastLeakMs = now
	if (l.water + 1) < l.capacity {
		// 尝试加水，并且水还未满
		l.water++
		return true
	}
	// 水满，拒绝加水
	return false
}

func (l *LeakyBucket) Set(r, c float64) {
	l.rate = r
	l.capacity = c
	l.water = 0
	l.lastLeakMs = time.Now().UnixNano() / 1e6
}
