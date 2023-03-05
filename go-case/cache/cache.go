package cache

import "time"

type Cache interface {
	// size : 1KB 100KB 1ME 2MB 1GB
	SetMaxMemory(size string) bool

	Set(key string, val any, expire time.Duration) bool

	Get(key string) (any, bool)

	Del(key string) bool

	Exists(key string) bool

	Flush() bool
	// 获取所有缓存中 key 的数量
	Keys() int64
}
