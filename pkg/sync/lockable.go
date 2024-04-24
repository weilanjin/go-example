package sync

import "sync"

type Lockable[T any] struct {
	mu    sync.Mutex
	value T
}

func (l *Lockable[T]) Get() T {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.value
}

func (l *Lockable[T]) Set(value T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.value = value
}
