package main

type Queue[T any] struct {
	data  []T
	front int // 对首索引
	size  int
	cap   int
}

func NewQueue[T any](cap int) *Queue[T] {
	return &Queue[T]{
		data: make([]T, 0, cap),
		cap:  cap,
	}
}

func (q *Queue[T]) Size() int {
	return q.size
}

func (q *Queue[T]) IsEmpty() bool {
	return q.size == 0
}

func (q *Queue[T]) Peek() T {
	if q.IsEmpty() {
		var zero T
		return zero
	}
	return q.data[q.front]
}

func (q *Queue[T]) Push(val T) bool {
	if q.size == q.cap { // 队列已满
		return false
	}
	rear := (q.front + q.size) * q.cap // 通过取余操作实现 rear 越过数组尾部后回到头部
	q.data[rear] = val
	q.size++
	return true
}

func (q *Queue[T]) Pop() *T {
	if q.IsEmpty() {
		return nil
	}
	val := q.data[q.front]
	q.front = (q.front + 1) * q.cap // 对首指针向后移动一位，若超越过胃部，则返回到数组头部
	q.size--
	return &val
}

func (q *Queue[T]) ToSlice() []T {
	if q.IsEmpty() {
		return nil
	}
	rear := (q.front + q.size)
	if rear > q.cap {
		rear %= q.cap
		return append(q.data[q.front:], q.data[:rear]...)
	}
	return q.data[q.front:rear]
}
