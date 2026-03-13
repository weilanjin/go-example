package main

import (
	"cmp"
	"iter"
)

// Slice collects stream elements into a slice.
// Slice 将流元素收集为切片。
func (s Stream[T]) Slice() []T {
	capHint := 0
	if s.sizeHint > 0 {
		capHint = s.sizeHint
	}
	items := make([]T, 0, capHint)
	s.seq(func(v T) bool {
		items = append(items, v)
		return true
	})
	return items
}

// All collects all stream elements into a slice.
// All 将全部流元素收集为切片。
func (s Stream[T]) All() []T { return s.Slice() }

// CollectToMap collects stream elements into a map by key/value functions.
// CollectToMap 按 key/value 函数将流元素收集为 map。
func CollectToMap[T any, K comparable, V any](s Stream[T], keyFn func(T) K, valueFn func(T) V) map[K]V {
	if keyFn == nil || valueFn == nil {
		return map[K]V{}
	}
	capHint := 0
	if s.sizeHint > 0 {
		capHint = s.sizeHint
	}
	result := make(map[K]V, capHint)
	s.seq(func(v T) bool {
		result[keyFn(v)] = valueFn(v)
		return true
	})
	return result
}

// ToMap is an alias of CollectToMap.
// ToMap 是 CollectToMap 的别名。
func ToMap[T any, K comparable, V any](s Stream[T], keyFn func(T) K, valueFn func(T) V) map[K]V {
	return CollectToMap(s, keyFn, valueFn)
}

// Len returns stream size, using size hint when available.
// Len 返回流长度，优先使用 size hint。
func (s Stream[T]) Len() int {
	if s.sizeHint >= 0 {
		return s.sizeHint
	}
	count := 0
	s.seq(func(T) bool {
		count++
		return true
	})
	return count
}

// Min returns minimal element by custom comparator.
// Min 按自定义比较器返回最小元素。
func (s Stream[T]) Min(less func(a, b T) bool) Optional[T] {
	if less == nil {
		return Optional[T]{}
	}
	first, ok := s.Head()
	if !ok {
		return Optional[T]{}
	}
	minV := first
	skippedFirst := false
	s.seq(func(v T) bool {
		if !skippedFirst {
			skippedFirst = true
			return true
		}
		if less(v, minV) {
			minV = v
		}
		return true
	})
	return Optional[T]{value: minV, present: true}
}

// Max returns maximal element by custom comparator.
// Max 按自定义比较器返回最大元素。
func (s Stream[T]) Max(less func(a, b T) bool) Optional[T] {
	if less == nil {
		return Optional[T]{}
	}
	first, ok := s.Head()
	if !ok {
		return Optional[T]{}
	}
	maxV := first
	skippedFirst := false
	s.seq(func(v T) bool {
		if !skippedFirst {
			skippedFirst = true
			return true
		}
		if less(maxV, v) {
			maxV = v
		}
		return true
	})
	return Optional[T]{value: maxV, present: true}
}

// MinOrdered returns minimal element for ordered types.
// MinOrdered 返回有序类型中的最小元素。
func MinOrdered[T cmp.Ordered](s Stream[T]) Optional[T] {
	return s.Min(func(a, b T) bool { return a < b })
}

// MaxOrdered returns maximal element for ordered types.
// MaxOrdered 返回有序类型中的最大元素。
func MaxOrdered[T cmp.Ordered](s Stream[T]) Optional[T] {
	return s.Max(func(a, b T) bool { return a < b })
}

// Any reports whether any element matches predicate.
// Any 判断是否存在满足条件的元素。
func (s Stream[T]) Any(predicate func(T) bool) bool {
	if predicate == nil {
		return false
	}
	matched := false
	s.seq(func(v T) bool {
		if predicate(v) {
			matched = true
			return false
		}
		return true
	})
	return matched
}

// Every reports whether all elements match predicate.
// Every 判断是否全部元素满足条件。
func (s Stream[T]) Every(predicate func(T) bool) bool {
	if predicate == nil {
		return false
	}
	all := true
	s.seq(func(v T) bool {
		if !predicate(v) {
			all = false
			return false
		}
		return true
	})
	return all
}

// None reports whether no elements match predicate.
// None 判断是否没有元素满足条件。
func (s Stream[T]) None(predicate func(T) bool) bool {
	return !s.Any(predicate)
}

// Head returns the first element and whether it exists.
// Head 返回首个元素及是否存在。
func (s Stream[T]) Head() (T, bool) {
	var first T
	found := false
	s.seq(func(v T) bool {
		first = v
		found = true
		return false
	})
	return first, found
}

// Each applies consumer to every element.
// Each 对每个元素执行 consumer。
func (s Stream[T]) Each(consumer func(T)) {
	if consumer == nil {
		return
	}
	s.seq(func(v T) bool {
		consumer(v)
		return true
	})
}

// EachUntil applies consumer until it returns false.
// EachUntil 执行到 consumer 返回 false 为止。
func (s Stream[T]) EachUntil(consumer func(T) bool) {
	if consumer == nil {
		return
	}
	s.seq(func(v T) bool {
		return consumer(v)
	})
}

// Seq exposes underlying iter.Seq for advanced integration.
// Seq 暴露底层 iter.Seq 以便高级集成。
func (s Stream[T]) Seq() iter.Seq[T] {
	return s.seq
}

// Iterator converts stream into a read-only channel iterator.
// Iterator 将流转换为只读 channel 迭代器。
func (s Stream[T]) Iterator() <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		s.seq(func(v T) bool {
			ch <- v
			return true
		})
	}()
	return ch
}

// Unordered marks stream as unordered.
// Unordered 将流标记为无序。
func (s Stream[T]) Unordered() Stream[T] {
	s.unordered = true
	return s
}

// Parallel marks stream as parallel-capable.
// Parallel 将流标记为可并行。
func (s Stream[T]) Parallel() Stream[T] {
	s.parallel = true
	return s
}

// Sequential marks stream as sequential mode.
// Sequential 将流标记为串行模式。
func (s Stream[T]) Sequential() Stream[T] {
	s.parallel = false
	return s
}

// IsParallel reports whether stream is marked parallel.
// IsParallel 返回流是否被标记为并行。
func (s Stream[T]) IsParallel() bool {
	return s.parallel
}

// OnClose registers close handlers executed in reverse order.
// OnClose 注册关闭回调，执行顺序为逆序。
func (s Stream[T]) OnClose(handler func()) Stream[T] {
	if handler == nil {
		return s
	}
	s.onClose = append(s.onClose, handler)
	return s
}

// Close executes registered close handlers in LIFO order.
// Close 按 LIFO 顺序执行已注册的关闭回调。
func (s Stream[T]) Close() {
	for i := len(s.onClose) - 1; i >= 0; i-- {
		s.onClose[i]()
	}
}
