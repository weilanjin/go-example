package main

import "iter"

/*
| 算法类型   | 推出时间   | 输出长度             | 哈希冲突  | 安全等级        | 应用                       |
|-----------|----------|---------------------|--------- |---------------|---------------------------|
| MD5       | 1992     | 128 bit             | 较多     | 低，已被成功攻击  | 已被弃用，仍用于数据完整性检查 |
| SHA-1     | 1995     | 160 bit             | 较多     | 低，已被成功攻击  | 已被弃用                   |
| SHA-2     | 2002     | 256/512 bit         | 很少     | 高              | 加密货币交易验证、数字签名等  |
| SHA-3     | 2008     | 224/256/384/512 bit | 很少     | 高              | 可用于替代 SHA-2           |
*/

// 链式地址哈希表
type HashMap struct {
	size        int       // 键值对数量
	cap         int       // 哈希表容量
	loadThres   float64   // 触发扩容的负载因子阈值
	extendRatio int       // 扩容倍数
	buckets     [][]*pair // 桶数组
}

type pair struct {
	key   int
	value string
}

func NewHashMap() *HashMap {
	buckets := make([][]*pair, 4)
	for i := range buckets {
		buckets[i] = make([]*pair, 0)
	}
	return &HashMap{
		size:        0,
		cap:         4,
		loadThres:   2.0 / 3.0,
		extendRatio: 2,
		buckets:     buckets,
	}
}

// 哈希函数
// 准确性、效率高、均匀分布
func (h *HashMap) hashFunc(key int) int {
	return key % h.cap
}

// 计算负载因子
func (h *HashMap) loadFactor() float64 {
	return float64(h.size) / float64(h.cap)
}

func (h *HashMap) Get(key int) string {
	index := h.hashFunc(key)
	bucket := h.buckets[index]
	// 遍历桶，若找到 key， 则返回对应 val
	for _, pair := range bucket {
		if pair.key == key {
			return pair.value
		}
	}
	return ""
}

func (h *HashMap) Put(key int, value string) {
	// 当负载因子超过阈值时，执行扩容
	if h.loadFactor() > h.loadThres {
		h.extend()
	}
	index := h.hashFunc(key)
	// 遍历桶，若找到 key， 则更新对应 val
	for i := range h.buckets[index] {
		if h.buckets[index][i].key == key {
			h.buckets[index][i].value = value
			return
		}
	}
	// 若没有找到 key， 则将 pair 添加到尾部
	pair := &pair{key, value}
	h.buckets[index] = append(h.buckets[index], pair)
	h.size++
}

func (h *HashMap) extend() {
	rawBuckets := make([][]*pair, len(h.buckets)) // 暂存原哈希表
	for i := range h.buckets {
		rawBuckets[i] = make([]*pair, len(h.buckets[i]))
		copy(rawBuckets[i], h.buckets[i])
	}
	// 扩展哈希表
	h.cap *= h.extendRatio
	h.buckets = make([][]*pair, h.cap)
	for i := range h.cap {
		h.buckets[i] = make([]*pair, 0)
	}
	h.size = 0
	// 将键值对从原哈希表搬运至新哈希表
	for _, bucket := range rawBuckets {
		for _, pair := range bucket {
			h.Put(pair.key, pair.value)
		}
	}
}

func (h *HashMap) Remove(key int) {
	index := h.hashFunc(key)
	for i, p := range h.buckets[index] {
		if p.key == key {
			// 切片删除
			h.buckets[index] = append(h.buckets[index][:i], h.buckets[index][i+1:]...)
			h.size--
			return
		}
	}
}

func (h *HashMap) Keys() iter.Seq[int] {
	return func(yield func(int) bool) {
		for _, bucket := range h.buckets {
			for _, pair := range bucket {
				if pair != nil {
					if !yield(pair.key) {
						return
					}
				}
			}
		}
	}
}

func (h *HashMap) Vals() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, bucket := range h.buckets {
			for _, pair := range bucket {
				if pair != nil {
					if !yield(pair.value) {
						return
					}
				}
			}
		}
	}
}

func (h *HashMap) All() iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		for _, bucket := range h.buckets {
			for _, pair := range bucket {
				if pair != nil {
					if !yield(pair.key, pair.value) {
						return
					}
				}
			}
		}
	}
}
