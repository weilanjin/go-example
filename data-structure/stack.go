package main

type Stack[T any] struct {
	data []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		data: make([]T, 0, 16),
	}
}

func (s *Stack[T]) Len() int {
	return len(s.data)
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.data) == 0
}

func (s *Stack[T]) Push(val T) {
	s.data = append(s.data, val)
}

func (s *Stack[T]) Pop() T {
	if len(s.data) == 0 {
		var zero T
		return zero
	}
	val := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return val
}

// Peek 返回栈顶元素，但不删除
func (s *Stack[T]) Peek() T {
	if s.IsEmpty() {
		var zero T
		return zero
	}
	return s.data[len(s.data)-1]
}

func (s *Stack[T]) ToSlice() []T {
	return s.data
}
