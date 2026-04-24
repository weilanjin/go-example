package sqlx

import (
	"iter"
)

type OrderedMap[K comparable, V any] struct {
	keys []K
	data map[K]V
}

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		keys: make([]K, 0),
		data: make(map[K]V),
	}
}

func (m *OrderedMap[K, V]) Set(key K, value V) {
	if _, exists := m.data[key]; !exists {
		m.keys = append(m.keys, key)
	}
	m.data[key] = value
}

func (m *OrderedMap[K, V]) Get(key K) (V, bool) {
	value, exists := m.data[key]
	return value, exists
}

func (m *OrderedMap[K, V]) Keys() []K {
	return m.keys
}

func (m *OrderedMap[K, V]) Values() []V {
	values := make([]V, len(m.keys))
	for i, key := range m.keys {
		values[i] = m.data[key]
	}
	return values
}

// All returns an iterator over all key-value pairs in the ordered map.
// All 返回有序映射中所有键值对的迭代器。
func (m *OrderedMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, key := range m.keys {
			if !yield(key, m.data[key]) {
				break
			}
		}
	}
}

func (m *OrderedMap[K, V]) Len() int {
	return len(m.keys)
}

func (m *OrderedMap[K, V]) Clear() {
	m.keys = m.keys[:0]
	m.data = make(map[K]V)
}

func (m *OrderedMap[K, V]) Has(key K) bool {
	_, exists := m.data[key]
	return exists
}

func (m *OrderedMap[K, V]) Delete(key K) {
	if _, exists := m.data[key]; exists {
		delete(m.data, key)
		for i, k := range m.keys {
			if k == key {
				m.keys = append(m.keys[:i], m.keys[i+1:]...)
				break
			}
		}
	}
}
