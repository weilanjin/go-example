// Package cache provides a cache implementation.
// 本地缓存的用处
// 1.高读写+命中率
// 2.减少网络请求
// 3.减少gc

// 实现高并发思路
// -> 数据分片(降低锁的粒度)

// 实现零GC方案
// 1.无GC
// -> 分配堆外内存(Mmap)
// 2.避免GC
// -> map 非指针优化(map[uint64]uint32)或者采用slice实现一套无指针的map
// —> 数据存入[]byte slice(可考虑底层采用环形队列封装循环使用空间)

// 1. freecache : https://github.com/coocood/freecache
// 2. bigcache : https://github.com/allegro/bigcache
// 3. fastcache : https://github.com/VictoriaMetrics/fastcache
// 5. groupcache : https://github.com/golang/groupcache
// 6. ristretto: https://github.com/dgraph-io/ristretto
// 7. go-cache : https://github.com/patrickmn/go-cache
package cache
