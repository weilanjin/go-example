package benchmark

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

// CommonMapWithRWMutex 结构体，包含普通 map 和读写锁
type CommonMapWithRWMutex struct {
	data    map[int]int
	rwMutex sync.RWMutex
}

// SyncMapStruct 结构体，包含 sync.Map
type SyncMapStruct struct {
	data sync.Map
}

// CommonMapWithRWMutex 的读方法
func (m *CommonMapWithRWMutex) Read(key int) int {
	// 获取读锁
	m.rwMutex.RLock()
	// 读取数据
	value, _ := m.data[key]
	// 释放读锁
	m.rwMutex.RUnlock()
	return value
}

// CommonMapWithRWMutex 的写方法
func (m *CommonMapWithRWMutex) Write(key, value int) {
	// 获取写锁
	m.rwMutex.Lock()
	// 写入数据
	m.data[key] = value
	// 释放写锁
	m.rwMutex.Unlock()
}

// SyncMapStruct 的读方法
func (s *SyncMapStruct) Read(key int) int {
	var value interface{}
	var ok bool
	// 从 sync.Map 中加载数据
	value, ok = s.data.Load(key)
	if ok {
		return value.(int)
	}
	return 0
}

// SyncMapStruct 的写方法
func (s *SyncMapStruct) Write(key, value int) {
	// 向 sync.Map 中存储数据
	s.data.Store(key, value)
}

// benchmarkCommonMapWithRWMutexRead 函数用于对 CommonMapWithRWMutex 的读操作进行压力测试
func benchmarkCommonMapWithRWMutexRead(b *testing.B, mapSize int) {
	commonMap := &CommonMapWithRWMutex{
		data: make(map[int]int),
	}
	// 初始化数据
	for i := 0; i < mapSize; i++ {
		commonMap.data[i] = i
	}

	var totalReadOps int64
	// 重置计时器
	b.ResetTimer()
	// 并发执行测试逻辑
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < mapSize; i++ {
				_ = commonMap.Read(i)
				// 原子操作增加读操作次数
				atomic.AddInt64(&totalReadOps, 1)
			}
		}
	})
	// 停止计时器
	b.StopTimer()
}

// benchmarkCommonMapWithRWMutexWrite 函数用于对 CommonMapWithRWMutex 的写操作进行压力测试
func benchmarkCommonMapWithRWMutexWrite(b *testing.B, mapSize int) {
	commonMap := &CommonMapWithRWMutex{
		data: make(map[int]int),
	}

	var totalWriteOps int64
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < mapSize; i++ {
				commonMap.Write(i, i)
				// 原子操作增加写操作次数
				atomic.AddInt64(&totalWriteOps, 1)
			}
		}
	})
	b.StopTimer()
}

// benchmarkSyncMapRead 函数用于对 SyncMap 的读操作进行压力测试
func benchmarkSyncMapRead(b *testing.B, mapSize int) {
	syncMap := &SyncMapStruct{}
	for i := 0; i < mapSize; i++ {
		syncMap.data.Store(i, i)
	}

	var totalReadOps int64
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < mapSize; i++ {
				_ = syncMap.Read(i)
				atomic.AddInt64(&totalReadOps, 1)
			}
		}
	})
	b.StopTimer()
}

// benchmarkSyncMapWrite 函数用于对 SyncMap 的写操作进行压力测试
func benchmarkSyncMapWrite(b *testing.B, mapSize int) {
	syncMap := &SyncMapStruct{}

	var totalWriteOps int64
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < mapSize; i++ {
				syncMap.Write(i, i)
				atomic.AddInt64(&totalWriteOps, 1)
			}
		}
	})
	b.StopTimer()
}

// TestPerformanceComparison 函数执行不同大小数据量下的压力测试
func Benchmark_RWMutex_Map(t *testing.B) {
	// sync.Map 和 RWMutex 的性能对比
	// 读多写少场景 sync.Map
	// 写多读少场景 RWMutex

	/*
		xxxxxxxxxxxxxxxxxxxxxxx mapSize 1000
		Benchmark_RWMutex_Map/RWMutexRead-1000-10         	    7777	    152460 ns/op
		Benchmark_RWMutex_Map/SyncMapRead-1000-10         	   15814	     76356 ns/op
		Benchmark_RWMutex_Map/RWMutexWrite-1000-10        	    4867	    212704 ns/op
		Benchmark_RWMutex_Map/SyncMapWrite-1000-10        	    3945	    275151 ns/op
		xxxxxxxxxxxxxxxxxxxxxxx mapSize 10000
		Benchmark_RWMutex_Map/RWMutexRead-10000-10        	     734	   1664277 ns/op
		Benchmark_RWMutex_Map/SyncMapRead-10000-10        	    1492	    804180 ns/op
		Benchmark_RWMutex_Map/RWMutexWrite-10000-10       	     685	   1808495 ns/op
		Benchmark_RWMutex_Map/SyncMapWrite-10000-10       	     374	   3158352 ns/op
		xxxxxxxxxxxxxxxxxxxxxxx mapSize 100000
		Benchmark_RWMutex_Map/RWMutexRead-100000-10       	     100	  16671079 ns/op
		Benchmark_RWMutex_Map/SyncMapRead-100000-10       	     146	   8087067 ns/op
		Benchmark_RWMutex_Map/RWMutexWrite-100000-10      	     100	  16988989 ns/op
		Benchmark_RWMutex_Map/SyncMapWrite-100000-10      	      34	  32041777 ns/op
	*/

	mapSizes := []int{1000, 10000, 100000}
	for _, mapSize := range mapSizes {
		fmt.Println("xxxxxxxxxxxxxxxxxxxxxxx", "mapSize", mapSize)

		// 对 CommonMapWithRWMutex 的读操作进行测试
		t.Run(fmt.Sprintf("RWMutexRead-%d", mapSize), func(t *testing.B) {
			benchmarkCommonMapWithRWMutexRead(t, mapSize)
		})
		// 对 SyncMap 的读操作进行测试
		t.Run(fmt.Sprintf("SyncMapRead-%d", mapSize), func(t *testing.B) {
			benchmarkSyncMapRead(t, mapSize)
		})

		// 对 CommonMapWithRWMutex 的写操作进行测试
		t.Run(fmt.Sprintf("RWMutexWrite-%d", mapSize), func(t *testing.B) {
			benchmarkCommonMapWithRWMutexWrite(t, mapSize)
		})
		// 对 SyncMap 的写操作进行测试
		t.Run(fmt.Sprintf("SyncMapWrite-%d", mapSize), func(t *testing.B) {
			benchmarkSyncMapWrite(t, mapSize)
		})
	}
}
