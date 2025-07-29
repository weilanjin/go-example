package main

import (
	"errors"
	"sync"
)

// Load Balancer Implementation 负载均衡器实现
type LeastConnections struct {
	pool *ServerPool
	mu   sync.Mutex
	idx  int // 当前服务器索引
}

func NewLeastConnectionsBalancer(pool *ServerPool) *LeastConnections {
	return &LeastConnections{
		pool: pool,
		idx:  -1, // 还没有开始选择任何服务器
	}
}

func (lc *LeastConnections) GetNextServer() (*Server, error) {
	if len(lc.pool.Servers) == 0 {
		return nil, errors.New("no servers found") // 没有可用服务器
	}
	lc.mu.Lock()
	defer lc.mu.Unlock()

	// 找到连接数最少的服务器
	minIndex := 0
	// minConnections := lc.pool.Servers[0].ConnectionCount

	// for i, server := range lc.pool.Servers {
	// 	if server.ConnectionCount < minConnections {
	// 		minConnections = server.ConnectionCount
	// 		minIndex = i
	// 	}
	// }

	return lc.pool.Servers[minIndex], nil
}
