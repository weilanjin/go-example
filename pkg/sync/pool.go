package sync

import "sync"

// 注意一下几点
// 1. sync.Pool 没有固定的大小.
// 2. 对象可能会被GC回收.
// 3. 对象有状态, 放入池子之前或之后,检索时清除或重置状态.

type Pool[T any] struct {
	internal *sync.Pool
}

func NewPool[T any](new func() T) Pool[T] {
	return Pool[T]{
		internal: &sync.Pool{
			New: func() any {
				return new()
			},
		},
	}
}

func (p Pool[T]) Get() T {
	return p.internal.Get().(T)
}

func (p Pool[T]) Put(t T) {
	p.internal.Put(t)
}
