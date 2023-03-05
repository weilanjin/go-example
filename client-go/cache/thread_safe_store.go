package cache

import "sync"

type ThreadSafeStore interface {
	Add(key string, obj any)
	Update(key string, obj any)
	Delete(key string)
	Get(key string) (item any, exists bool)
	List() []any
	ListKeys() []string
	Replace(map[string]any, string)
	Index(indexName string, obj any) ([]any, error)
	IndexKeys(indexName, indexedValue string) ([]string, error)
	ListIndexFuncValues(name string) []string
	ByIndex(indexName, indexedValue string) ([]any, error)
	GetIndexers() Indexers
	AddIndexers(newIndexers Indexers) error
	Resync() error
}

type threadSafeMap struct {
	lock  sync.RWMutex
	items map[string]any

	index *storeIndex
}

func NewThreadSafeStore(indexers Indexers, indices Indices) ThreadSafeStore {
	// return &threadSafeMap{
	// 	items: map[string]any{},
	// 	index: &storeIndex{
	// 		indexers: indexers,
	// 		indices:  indices,
	// 	},
	// }
	return nil
}

func (c *threadSafeMap) Add(key string, obj any) {
	c.Update(key, obj)
}

func (c *threadSafeMap) Update(key string, obj any) {
	c.lock.Lock()
	defer c.lock.Unlock()

	oldItem := c.items[key]
	c.items[key] = obj
	c.index.updateIndices(oldItem, obj, key)
}

func (c *threadSafeMap) Delete(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if obj, ok := c.items[key]; ok {
		c.index.updateIndices(obj, nil, key)
		delete(c.items, key)
	}
}

func (c *threadSafeMap) Resync() error {
	return nil
}

type storeIndex struct {
	indexers Indexers
	indices  Indices
}

func (i *storeIndex) updateIndices(oldObj, newObj any, key string) {}
