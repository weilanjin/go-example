package mutex

import "sync"

// carlmjohnson/syncx

type Value[T any] struct {
	mu    sync.Mutex
	value T
}

func NewValue[T any](initial T) *Value[T] {
	var m Value[T]
	m.value = initial
	return &m
}

func (m *Value[T]) Lock(f func(value *T)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v := m.value
	f(&v)
	m.value = v
}

func (m *Value[T]) Load() T {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.value
}

func (m *Value[T]) Store(value T) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.value = value
}
