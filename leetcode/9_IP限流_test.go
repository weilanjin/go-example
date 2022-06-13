package leetcode

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 场景：在⼀个⾼并发的web服务器中，要限制IP的频繁访问。现模拟100个IP同时并发访问服
// 务器，每个IP要重复访问1000次。
// 每个IP三分钟之内只能访问⼀次

type Ban struct {
	visitIPs map[string]time.Time
	lock     sync.Mutex
}

func NewBan() *Ban {
	return &Ban{visitIPs: make(map[string]time.Time)}
}

func (o *Ban) visit(ip string) bool {
	o.lock.Lock()
	defer o.lock.Unlock()
	if t, ok := o.visitIPs[ip]; ok {
		// t + 3 是否在 now 之后
		// true 还没过期
		// false 已过期
		after := t.Add(time.Minute * 3).After(time.Now())
		if after {
			return after
		}
	}
	o.visitIPs[ip] = time.Now()
	return false
}

func Test9(t *testing.T) {
	success := int64(0) // 原子性
	ban := NewBan()
	for i := 0; i < 1000; i++ {
		for j := 0; j < 100; j++ {
			go func(j int) {
				ip := fmt.Sprintf("192.168.1.%d", j)
				if !ban.visit(ip) {
					atomic.AddInt64(&success, 1)
				}
			}(j)
		}
	}
	log.Println("success:", success)
}
