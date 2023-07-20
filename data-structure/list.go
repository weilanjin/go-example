package main

// List doubly-linked list
// container/list/list.go
type List[T any] struct {
	root *Node[T]
	len  int
}

func NewList[T any]() *List[T] {
	return &List[T]{}
}

func (l *List[T]) Len() int {
	return l.len
}

func (l *List[T]) PushFront(v T) {
	n := Node[T]{
		Value: v,
		Next:  l.root.Prev,
	}
	l.root.Next = &n
	l.len++
}

func (l *List[T]) PushBack(v T) {

}