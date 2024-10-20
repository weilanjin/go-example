package queue

import (
	"sync/atomic"
	"unsafe"
)

// Maged M. Micheal 和 Micheal L. Scott 1996 年发布的论文 "Simple, Fast, and Practical Non-Blocking and Blocking Concurrent Queue Algorithms"

// LKQueue 是以 lock-free 方式实现的队列, 它只需要head和tail两个字段
type LKQueue[T any] struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

// 队列中的每个节点, 除自己的值外, 还有 next 字段指向下一个节点
type node[T any] struct {
	value T
	next  unsafe.Pointer
}

func NewLKQueue[T any]() *LKQueue[T] {
	n := unsafe.Pointer(&node[T]{})
	return &LKQueue[T]{head: n, tail: n}
}

// 入队
func (q *LKQueue[T]) Enqueue(value T) {
	n := &node[T]{value: value}
	for {
		tail := load[T](&q.tail)
		next := load[T](&tail.next)
		if tail == load[T](&q.tail) { // tail 和 next 是否一致
			if next == nil {
				if cas(&tail.next, next, n) {
					cas(&q.tail, tail, n) // 入队完成, 设置tail
					return
				}
			} else {
				// 队列中间节点
				cas(&q.tail, tail, next)
			}
		}
	}
}

// 出队
func (q *LKQueue[T]) Dequeue() T {
	var t T
	for {
		head := load[T](&q.head)
		tail := load[T](&q.tail)
		next := load[T](&head.next)
		if head == load[T](&q.head) { // head 和 tail 是否一致
			if head == tail { // 队列为空, 或者tail还未到队尾
				if next == nil { // 队列为空
					return t
				}
				// 队列中有元素, tail 指向最后一个节点
				cas(&q.tail, tail, next)
			} else {
				// 队列中有元素, 取出第一个节点
				v := next.value
				if cas(&q.head, head, next) {
					return v // 出队完成
				}
			}
		}
	}
}

// 读取节点的值
func load[T any](p *unsafe.Pointer) (n *node[T]) {
	return (*node[T])(atomic.LoadPointer(p))
}

// 原子地修改节点的值
func cas[T any](p *unsafe.Pointer, old, new *node[T]) bool {
	return atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(p)), unsafe.Pointer(old), unsafe.Pointer(new))
}
