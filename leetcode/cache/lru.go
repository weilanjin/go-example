package cache

import (
	"container/list"
)

// LRUCache https://github.com/kubernetes/utils/blob/master/lru/lru.go
// Least recently used，最近最少使用
// 哈希表 + 双向链表
// 1. LRU 双端链表实现，访问到的节点移动到头部，超出容量的从尾部删除
// 2. 实现O(1)使用map
type LRUCache struct {
	cap   int
	cache map[string]*list.Element
	ll    *list.List
}

type entry struct {
	key   string
	value any
}

func New(cap int) LRUCache {
	return LRUCache{
		cap:   cap,
		cache: make(map[string]*list.Element),
		ll:    list.New(),
	}
}

// GET
// 1. 直接从map获取
// 2. 把当前访问的元数移动到头部
func (c *LRUCache) GET(key string) (any, bool) {
	if c.cache == nil {
		return nil, false
	}
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToFront(e)
		return e.Value.(*entry).value, true
	}
	return nil, false
}

// PUT
// 1.校验 map linked list 是否 clear
// 2.是否已经存在，存在则修改值，链表元数移至头部
// 3.不存在则push到头部.
// 4.如果容量超过，就驱除尾部一个元素
func (c *LRUCache) PUT(key string, value any) {
	if c.cache == nil { // clear
		c.cache = make(map[string]*list.Element)
		c.ll = list.New()
	}
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToFront(e)
		e.Value.(*entry).value = value
		return
	}
	e := c.ll.PushFront(&entry{key, value})
	c.cache[key] = e
	if c.cap != 0 && c.ll.Len() > c.cap {
		c.removeOldest()
	}
}

// Remove
// 同时移除 map linkedlist
func (c *LRUCache) Remove(key string) {
	if c.cache == nil {
		return
	}
	if e, ok := c.cache[key]; ok {
		c.removeElement(e)
	}
}

func (c *LRUCache) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

func (c *LRUCache) Clear() {
	c.ll = nil
	c.cache = nil
}

func (c *LRUCache) removeOldest() {
	if c.cache == nil {
		return
	}
	if e := c.ll.Back(); e != nil {
		c.removeElement(e)
	}
}

func (c *LRUCache) removeElement(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cache, kv.key)
}