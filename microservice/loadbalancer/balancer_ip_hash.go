package main

import (
	"errors"
	"sync"
)

type IpHash struct {
	pool *ServerPool
	mu   sync.Mutex
	idx  int // 当前服务器索引
}

func NewIpHashBalancer(pool *ServerPool) *IpHash {
	return &IpHash{
		pool: pool,
		idx:  -1, // 还没有开始选择任何服务器
	}
}

func (ih *IpHash) GetNextServer(ip string) (*Server, error) {
	if len(ih.pool.Servers) == 0 {
		return nil, errors.New("no servers found") // 没有可用服务器
	}
	ih.mu.Lock()
	defer ih.mu.Unlock()

	// 使用 IP 哈希算法选择服务器
	hash := hashIP(ip)
	serverIndex := hash % len(ih.pool.Servers)

	return ih.pool.Servers[serverIndex], nil
}

// hashIP 是一个简单的哈希函数，用于将 IP 地址转换为整数
func hashIP(ip string) int {
	hash := 0
	for i := 0; i < len(ip); i++ {
		hash = (hash << 5) - hash + int(ip[i]) // 简单的哈希算法
	}
	return hash
}
