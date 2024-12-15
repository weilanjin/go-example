package mutex

import "sync"

type Slice[T any] struct {
	mu sync.Mutex
	s  []T
}

func NewSlice[T any](size ...int) *Slice[T] {
	if len(size) > 0 {
		return &Slice[T]{
			s: make([]T, size[0]),
		}
	}
	return &Slice[T]{
		s: make([]T, 0),
	}
}

func (s *Slice[T]) Get(index int) (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if index < 0 || index >= len(s.s) {
		return *new(T), false
	}
	return s.s[index], true
}

func (s *Slice[T]) Set(index int, value T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if index < 0 || index >= len(s.s) {
		return
	}
	s.s[index] = value
}

func (s *Slice[T]) Append(value T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.s = append(s.s, value)
}

func (s *Slice[T]) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.s)
}

func (s *Slice[T]) Cap() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return cap(s.s)
}

func (s *Slice[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.s = s.s[:0]
}

func (s *Slice[T]) Slice(start, end int) []T {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.s[start:end]
}

func (s *Slice[T]) SliceFrom(start int) []T {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.s[start:]
}

func (s *Slice[T]) SliceTo(end int) []T {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.s[:end]
}

func (s *Slice[T]) Swap(i, j int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.s[i], s.s[j] = s.s[j], s.s[i]
}

func (s *Slice[T]) Reverse() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, j := 0, len(s.s)-1; i < j; i, j = i+1, j-1 {
		s.s[i], s.s[j] = s.s[j], s.s[i]
	}
}
