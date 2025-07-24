package main

import (
	"errors"
	"sync"
)

// round-robin 轮询算法
type RoundRobinBalancer struct {
	pool *ServerPool
	mu   sync.Mutex
	idx  int // 当前服务器索引
}

func NewRoundRobinBalancer(pool *ServerPool) *RoundRobinBalancer {
	return &RoundRobinBalancer{
		pool: pool,
		idx:  -1, // 还没有开始选择任何服务器
	}
}

func (rr *RoundRobinBalancer) GetNextServer() (*Server, error) {
	if len(rr.pool.Servers) == 0 {
		return nil, errors.New("no servers found") // 没有可用服务器
	}
	rr.mu.Lock()
	defer rr.mu.Unlock()

	// 增加索引并循环到第一个服务器
	rr.idx = (rr.idx + 1) % len(rr.pool.Servers)
	return rr.pool.Servers[rr.idx], nil
}
