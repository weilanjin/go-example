package main

type Deque[T any] struct {
	data  []T // 用于储存双向队列元素的数组
	front int // 队首指针，指向队首元素
	size  int
	cap   int
}

func NewDeque[T any](cap int) *Deque[T] {
	return &Deque[T]{
		data: make([]T, cap),
		cap:  cap,
	}
}

func (d *Deque[T]) Size() int {
	return d.size
}

func (d *Deque[T]) IsEmpty() bool {
	return d.size == 0
}

func (d *Deque[T]) index(i int) int {
	// 通过取余操作实现数组首位相连
	// 当 i 超过数组尾部后，回到头部
	// 当 i 超越数组头部时，回到尾部
	return (i + d.cap) % d.cap
}

func (d *Deque[T]) PushFirst(e T) bool {
	if d.size == d.cap {
		return false
	}
	d.front = d.index(d.front - 1)
	d.data[d.front] = e
	d.size++
	return true
}

func (d *Deque[T]) PushLast(e T) bool {
	if d.size == d.cap {
		return false
	}
	tail := d.index(d.front + d.size)
	d.data[tail] = e
	d.size++
	return true
}

func (d *Deque[T]) PopFirst() *T {
	if d.size == 0 {
		return nil
	}
	e := d.data[d.front]
	d.front = d.index(d.front + 1) // 队首指针后移一位
	d.size--
	return &e
}

func (d *Deque[T]) PopLast() *T {
	e := d.PeekLast()
	if e == nil {
		return nil
	}
	d.size--
	return e
}

func (d *Deque[T]) PeekFirst() *T {
	if d.IsEmpty() {
		return nil
	}
	return &d.data[d.front]
}

func (d *Deque[T]) PeekLast() *T {
	if d.IsEmpty() {
		return nil
	}
	tail := d.index(d.front + d.size - 1)
	return &d.data[tail]
}

func (d *Deque[T]) ToSlice() []T {
	// 仅转换有效长度范围内的列表元素
	res := make([]T, d.size)
	for i, j := 0, d.front; i < d.size; i++ {
		res[i] = d.data[d.index(j)]
		j++
	}
	return res
}
