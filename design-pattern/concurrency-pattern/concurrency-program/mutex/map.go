package mutex

import "sync"

type Map[k comparable, v any] struct {
	mu sync.Mutex
	m  map[k]v
}

func NewMap[k comparable, v any](size ...int) *Map[k, v] {
	if len(size) > 0 {
		return &Map[k, v]{
			m: make(map[k]v, size[0]),
		}
	}
	return &Map[k, v]{
		m: make(map[k]v),
	}
}

func (m *Map[k, v]) Get(key k) (v, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	value, ok := m.m[key]
	return value, ok
}

func (m *Map[k, v]) Set(key k, value v) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m[key] = value
}

func (m *Map[k, v]) Delete(key k) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.m, key)
}

func (m *Map[k, v]) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.m)
}

func (m *Map[k, v]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m = make(map[k]v)
}

func (m *Map[k, v]) Range(f func(key k, value v) bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for key, value := range m.m {
		if !f(key, value) {
			break
		}
	}
}

func (m *Map[k, v]) Keys() []k {
	m.mu.Lock()
	defer m.mu.Unlock()
	keys := make([]k, 0, len(m.m))
	for key := range m.m {
		keys = append(keys, key)
	}
	return keys
}

func (m *Map[k, v]) Values() []v {
	m.mu.Lock()
	defer m.mu.Unlock()
	values := make([]v, 0, len(m.m))
	for _, value := range m.m {
		values = append(values, value)
	}
	return values
}

func (m *Map[k, v]) Copy() *Map[k, v] {
	m.mu.Lock()
	defer m.mu.Unlock()
	newMap := NewMap[k, v](len(m.m))
	for key, value := range m.m {
		newMap.m[key] = value
	}
	return newMap
}

func (m *Map[k, v]) Clone() *Map[k, v] {
	m.mu.Lock()
	defer m.mu.Unlock()
	newMap := NewMap[k, v](len(m.m))
	for key, value := range m.m {
		newMap.m[key] = value
	}
	return newMap
}
