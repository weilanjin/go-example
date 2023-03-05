// 内存相关工具包
package mem

import "runtime"

// gc 回收后内存占用 byte
func MemConsumed() uint64 {
	runtime.GC()
	var s runtime.MemStats
	runtime.ReadMemStats(&s)
	return s.Sys
}
