package cache

import (
	"regexp"
	"strings"
)

type memCache struct {
	maxMemorySize int64
}

// size : 1KB 100KB 1ME 2MB 1GB
func (mc *memCache) SetMaxMemory(size string) bool {

	re, _ := regexp.Compile("[0-9]+")
	// 去掉数字部分
	unit := string(re.ReplaceAllString(size, ""))
	unit = strings.ToUpper(unit)
	switch unit {
	case "B":
	case "KB":
	case "MB":
	case "GB":
	case "TB":
	case "PB":
	default:

	}
	return false
}

func (mc *memCache) Set(key string, val any) bool {
	return false
}

func (mc *memCache) Get(key string) (interface{}, bool) {
	return nil, false
}
func (mc *memCache) Del(key string) bool {
	return false
}

func (mc *memCache) Exists(key string) bool {
	return false
}

func (mc *memCache) Flush() bool {
	return false
}

func (mc *memCache) Keys() int64 {
	return 0
}
