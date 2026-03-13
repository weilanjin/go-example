package main

import "cmp"

func (s Stream[T]) CollectToSlice() []T {
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

func (s Stream[T]) ToArray() []T { return s.CollectToSlice() }
func (s Stream[T]) ToList() []T  { return s.CollectToSlice() }

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

func (s Stream[T]) Count() int {
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

func (s Stream[T]) Min(less func(a, b T) bool) Optional[T] {
	if less == nil {
		return Optional[T]{}
	}
	first, ok := s.FindFirst()
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

func (s Stream[T]) Max(less func(a, b T) bool) Optional[T] {
	if less == nil {
		return Optional[T]{}
	}
	first, ok := s.FindFirst()
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

func MinOrdered[T cmp.Ordered](s Stream[T]) Optional[T] {
	return s.Min(func(a, b T) bool { return a < b })
}

func MaxOrdered[T cmp.Ordered](s Stream[T]) Optional[T] {
	return s.Max(func(a, b T) bool { return a < b })
}

func (s Stream[T]) AnyMatch(predicate func(T) bool) bool {
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

func (s Stream[T]) AllMatch(predicate func(T) bool) bool {
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

func (s Stream[T]) NoneMatch(predicate func(T) bool) bool {
	return !s.AnyMatch(predicate)
}

func (s Stream[T]) FindFirst() (T, bool) {
	var first T
	found := false
	s.seq(func(v T) bool {
		first = v
		found = true
		return false
	})
	return first, found
}

func (s Stream[T]) FindAny() (T, bool) {
	return s.FindFirst()
}

func (s Stream[T]) ForEach(consumer func(T) bool) {
	if consumer == nil {
		return
	}
	s.seq(func(v T) bool {
		return consumer(v)
	})
}

func (s Stream[T]) ForEachOrdered(consumer func(T) bool) {
	s.ForEach(consumer)
}

func (s Stream[T]) Seq() Seq[T] {
	return s.seq
}

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

func (s Stream[T]) Unordered() Stream[T] {
	s.unordered = true
	return s
}

func (s Stream[T]) Parallel() Stream[T] {
	s.parallel = true
	return s
}

func (s Stream[T]) Sequential() Stream[T] {
	s.parallel = false
	return s
}

func (s Stream[T]) IsParallel() bool {
	return s.parallel
}

func (s Stream[T]) OnClose(handler func()) Stream[T] {
	if handler == nil {
		return s
	}
	s.onClose = append(s.onClose, handler)
	return s
}

func (s Stream[T]) Close() {
	for i := len(s.onClose) - 1; i >= 0; i-- {
		s.onClose[i]()
	}
}
