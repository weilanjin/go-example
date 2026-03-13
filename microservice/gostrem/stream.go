package main

import (
	"cmp"
	"iter"
	"sort"
)

// Stream represents a lazy pipeline over Seq.
// Stream 表示基于 Seq 的惰性处理流水线。
type Stream[T any] struct {
	seq        iter.Seq[T]
	sizeHint   int
	parallel   bool
	workerPool int
	unordered  bool
	onClose    []func()
}

// Tuple2 models Spark-style tuple2.
// Tuple2 表示 Spark 风格二元组。
type Tuple2[A any, B any] struct {
	First  A
	Second B
}

// Pair models key-value records for Spark-style pair operators.
// Pair 表示 Spark 风格 key-value 记录。
type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

func makeStream[T any](seq iter.Seq[T], sizeHint int) Stream[T] {
	if seq == nil {
		return Stream[T]{seq: func(func(T) bool) {}, sizeHint: 0}
	}
	if sizeHint < -1 {
		sizeHint = -1
	}
	return Stream[T]{seq: seq, sizeHint: sizeHint}
}

func (s Stream[T]) derive(seq iter.Seq[T], sizeHint int) Stream[T] {
	next := makeStream(seq, sizeHint)
	next.parallel = s.parallel
	next.workerPool = s.workerPool
	next.unordered = s.unordered
	if len(s.onClose) > 0 {
		next.onClose = append([]func(){}, s.onClose...)
	}
	return next
}

// Empty returns an empty stream.
// Empty 返回空流。
func Empty[T any]() Stream[T] {
	return makeStream[T](nil, 0)
}

// From creates a stream from values.
// From 从可变参数创建流。
func From[T any](values ...T) Stream[T] {
	return FromSlice(values)
}

// FromPointer creates empty stream when pointer is nil.
// FromPointer 在指针为 nil 时返回空流。
func FromPointer[T any](value *T) Stream[T] {
	if value == nil {
		return Empty[T]()
	}
	return From(*value)
}

// FromSlice creates a stream from a slice lazily.
// FromSlice 基于切片惰性创建流。
func FromSlice[T any](values []T) Stream[T] {
	if len(values) == 0 {
		return makeStream[T](nil, 0)
	}
	seq := func(yield func(T) bool) {
		for i := range values {
			if !yield(values[i]) {
				return
			}
		}
	}
	return makeStream(seq, len(values))
}

// Concat concatenates two streams lazily.
// Concat 惰性拼接两个流。
func Concat[T any](left, right Stream[T]) Stream[T] {
	seq := func(yield func(T) bool) {
		if !left.forEachUntil(yield) {
			return
		}
		right.forEachUntil(yield)
	}
	hint := -1
	if left.sizeHint >= 0 && right.sizeHint >= 0 {
		hint = left.sizeHint + right.sizeHint
	}
	res := makeStream(seq, hint)
	res.parallel = left.parallel || right.parallel
	if left.workerPool > right.workerPool {
		res.workerPool = left.workerPool
	} else {
		res.workerPool = right.workerPool
	}
	res.unordered = left.unordered || right.unordered
	if len(left.onClose) > 0 {
		res.onClose = append(res.onClose, left.onClose...)
	}
	if len(right.onClose) > 0 {
		res.onClose = append(res.onClose, right.onClose...)
	}
	return res
}

// Chain concatenates multiple streams lazily.
// Chain 惰性拼接多个流。
func Chain[T any](streams ...Stream[T]) Stream[T] {
	if len(streams) == 0 {
		return Empty[T]()
	}
	res := streams[0]
	for i := 1; i < len(streams); i++ {
		res = Concat(res, streams[i])
	}
	return res
}

// Generate creates an infinite stream from supplier.
// Generate 基于 supplier 创建无限流。
func Generate[T any](supplier func() T) Stream[T] {
	if supplier == nil {
		return Empty[T]()
	}
	seq := func(yield func(T) bool) {
		for {
			if !yield(supplier()) {
				return
			}
		}
	}
	return makeStream(seq, -1)
}

// Iterate creates an infinite iterative stream.
// Iterate 创建无限迭代流。
func Iterate[T any](seed T, next func(T) T) Stream[T] {
	if next == nil {
		return From(seed)
	}
	seq := func(yield func(T) bool) {
		v := seed
		for {
			if !yield(v) {
				return
			}
			v = next(v)
		}
	}
	return makeStream(seq, -1)
}

// IterateWhile creates finite iterative stream with hasNext check.
// IterateWhile 基于 hasNext 创建有限迭代流。
func IterateWhile[T any](seed T, hasNext func(T) bool, next func(T) T) Stream[T] {
	if hasNext == nil || next == nil {
		return Empty[T]()
	}
	seq := func(yield func(T) bool) {
		for v := seed; hasNext(v); v = next(v) {
			if !yield(v) {
				return
			}
		}
	}
	return makeStream(seq, -1)
}

// Builder accumulates values then builds a stream.
// Builder 先累积元素再构建流。
type Builder[T any] struct {
	values []T
}

// NewBuilder creates a stream builder.
// NewBuilder 创建流构建器。
func NewBuilder[T any]() *Builder[T] {
	return &Builder[T]{values: make([]T, 0)}
}

// Add appends one element into builder.
// Add 向构建器追加一个元素。
func (b *Builder[T]) Add(value T) *Builder[T] {
	b.values = append(b.values, value)
	return b
}

// Build creates stream from builder values.
// Build 从构建器元素创建流。
func (b *Builder[T]) Build() Stream[T] {
	return FromSlice(b.values)
}

func (s Stream[T]) forEachUntil(fn func(T) bool) bool {
	stopped := false
	s.seq(func(v T) bool {
		if !fn(v) {
			stopped = true
			return false
		}
		return true
	})
	return !stopped
}

// Filter keeps values that match predicate lazily.
// Filter 惰性保留满足条件的元素。
func (s Stream[T]) Filter(predicate func(T) bool) Stream[T] {
	if predicate == nil {
		return s
	}
	seq := func(yield func(T) bool) {
		s.seq(func(v T) bool {
			if predicate(v) {
				return yield(v)
			}
			return true
		})
	}
	return s.derive(seq, -1)
}

// Where is an alias of Filter.
// Where 是 Filter 的别名。
func (s Stream[T]) Where(predicate func(T) bool) Stream[T] {
	return s.Filter(predicate)
}

// Peek executes action for each value while preserving stream.
// Peek 在不改变流的情况下观察元素。
func (s Stream[T]) Peek(action func(T)) Stream[T] {
	if action == nil {
		return s
	}
	seq := func(yield func(T) bool) {
		s.seq(func(v T) bool {
			action(v)
			return yield(v)
		})
	}
	return s.derive(seq, s.sizeHint)
}

// Limit keeps at most n elements lazily.
// Limit 惰性保留前 n 个元素。
func (s Stream[T]) Limit(n int) Stream[T] {
	if n <= 0 {
		return makeStream[T](nil, 0)
	}
	seq := func(yield func(T) bool) {
		count := 0
		s.seq(func(v T) bool {
			if count >= n {
				return false
			}
			count++
			return yield(v)
		})
	}
	hint := -1
	if s.sizeHint >= 0 {
		if s.sizeHint < n {
			hint = s.sizeHint
		} else {
			hint = n
		}
	}
	return s.derive(seq, hint)
}

// Skip ignores the first n elements lazily.
// Skip 惰性跳过前 n 个元素。
func (s Stream[T]) Skip(n int) Stream[T] {
	if n <= 0 {
		return s
	}
	seq := func(yield func(T) bool) {
		skipped := 0
		s.seq(func(v T) bool {
			if skipped < n {
				skipped++
				return true
			}
			return yield(v)
		})
	}
	hint := -1
	if s.sizeHint >= 0 {
		if s.sizeHint <= n {
			hint = 0
		} else {
			hint = s.sizeHint - n
		}
	}
	return s.derive(seq, hint)
}

// TakeWhile keeps prefix while predicate is true.
// TakeWhile 在条件为真时保留前缀元素。
func (s Stream[T]) TakeWhile(predicate func(T) bool) Stream[T] {
	if predicate == nil {
		return Empty[T]()
	}
	seq := func(yield func(T) bool) {
		s.seq(func(v T) bool {
			if !predicate(v) {
				return false
			}
			return yield(v)
		})
	}
	return s.derive(seq, -1)
}

// DropWhile drops prefix while predicate is true.
// DropWhile 在条件为真时丢弃前缀元素。
func (s Stream[T]) DropWhile(predicate func(T) bool) Stream[T] {
	if predicate == nil {
		return s
	}
	seq := func(yield func(T) bool) {
		dropping := true
		s.seq(func(v T) bool {
			if dropping && predicate(v) {
				return true
			}
			dropping = false
			return yield(v)
		})
	}
	return s.derive(seq, -1)
}

// Distinct removes duplicate values lazily.
// Distinct 惰性去重。
func Distinct[T comparable](s Stream[T]) Stream[T] {
	seq := func(yield func(T) bool) {
		seen := make(map[T]struct{})
		s.seq(func(v T) bool {
			if _, ok := seen[v]; ok {
				return true
			}
			seen[v] = struct{}{}
			return yield(v)
		})
	}
	return s.derive(seq, -1)
}

// Unique is an alias of Distinct.
// Unique 是 Distinct 的别名。
func Unique[T comparable](s Stream[T]) Stream[T] {
	return Distinct(s)
}

// DistinctBy removes duplicates by comparable key lazily.
// DistinctBy 按可比较 key 惰性去重。
func DistinctBy[T any, K comparable](s Stream[T], keyFn func(T) K) Stream[T] {
	if keyFn == nil {
		return s
	}
	seq := func(yield func(T) bool) {
		seen := make(map[K]struct{})
		s.seq(func(v T) bool {
			k := keyFn(v)
			if _, ok := seen[k]; ok {
				return true
			}
			seen[k] = struct{}{}
			return yield(v)
		})
	}
	return s.derive(seq, -1)
}

// UniqueBy is an alias of DistinctBy.
// UniqueBy 是 DistinctBy 的别名。
func UniqueBy[T any, K comparable](s Stream[T], keyFn func(T) K) Stream[T] {
	return DistinctBy(s, keyFn)
}

// Sorted materializes then sorts values with less comparator.
// Sorted 会先收集再按 less 排序。
func (s Stream[T]) Sorted(less func(a, b T) bool) Stream[T] {
	if less == nil {
		return s
	}
	seq := func(yield func(T) bool) {
		items := s.Slice()
		sort.Slice(items, func(i, j int) bool {
			return less(items[i], items[j])
		})
		for i := range items {
			if !yield(items[i]) {
				return
			}
		}
	}
	return s.derive(seq, s.sizeHint)
}

// SortedByKey sorts by comparable key eagerly.
// SortedByKey 按可比较 key 进行排序（会收集后排序）。
func SortedByKey[T any, K cmp.Ordered](s Stream[T], keyFn func(T) K) Stream[T] {
	if keyFn == nil {
		return s
	}
	return s.Sorted(func(a, b T) bool {
		return keyFn(a) < keyFn(b)
	})
}

// Map transforms stream values lazily.
// Map 惰性转换流元素。
func Map[T any, R any](s Stream[T], mapper func(T) R) Stream[R] {
	if mapper == nil {
		return Empty[R]()
	}
	seq := func(yield func(R) bool) {
		s.seq(func(v T) bool {
			return yield(mapper(v))
		})
	}
	res := makeStream(seq, s.sizeHint)
	res.parallel = s.parallel
	res.unordered = s.unordered
	if len(s.onClose) > 0 {
		res.onClose = append([]func(){}, s.onClose...)
	}
	return res
}

// MapFn is an alias of Map.
// MapFn 是 Map 的别名。
func MapFn[T any, R any](s Stream[T], mapper func(T) R) Stream[R] {
	return Map(s, mapper)
}

// FlatMap expands each element to another stream lazily.
// FlatMap 惰性展开每个元素为子流。
func FlatMap[T any, R any](s Stream[T], mapper func(T) Stream[R]) Stream[R] {
	if mapper == nil {
		return Empty[R]()
	}
	seq := func(yield func(R) bool) {
		stop := false
		s.seq(func(v T) bool {
			if stop {
				return false
			}
			next := mapper(v)
			next.seq(func(r R) bool {
				if !yield(r) {
					stop = true
					return false
				}
				return true
			})
			return !stop
		})
	}
	res := makeStream(seq, -1)
	res.parallel = s.parallel
	res.unordered = s.unordered
	if len(s.onClose) > 0 {
		res.onClose = append([]func(){}, s.onClose...)
	}
	return res
}

// FlatMapFn is an alias of FlatMap.
// FlatMapFn 是 FlatMap 的别名。
func FlatMapFn[T any, R any](s Stream[T], mapper func(T) Stream[R]) Stream[R] {
	return FlatMap(s, mapper)
}

// MapToInt maps to int stream.
// MapToInt 映射为 int 流。
func MapToInt[T any](s Stream[T], mapper func(T) int) Stream[int] {
	return Map(s, mapper)
}

// MapToInt64 maps to int64 stream.
// MapToInt64 映射为 int64 流。
func MapToInt64[T any](s Stream[T], mapper func(T) int64) Stream[int64] {
	return Map(s, mapper)
}

// MapToFloat64 maps to float64 stream.
// MapToFloat64 映射为 float64 流。
func MapToFloat64[T any](s Stream[T], mapper func(T) float64) Stream[float64] {
	return Map(s, mapper)
}

// FlatMapToInt flat maps to int stream.
// FlatMapToInt 扁平映射为 int 流。
func FlatMapToInt[T any](s Stream[T], mapper func(T) Stream[int]) Stream[int] {
	return FlatMap(s, mapper)
}

// FlatMapToInt64 flat maps to int64 stream.
// FlatMapToInt64 扁平映射为 int64 流。
func FlatMapToInt64[T any](s Stream[T], mapper func(T) Stream[int64]) Stream[int64] {
	return FlatMap(s, mapper)
}

// FlatMapToFloat64 flat maps to float64 stream.
// FlatMapToFloat64 扁平映射为 float64 流。
func FlatMapToFloat64[T any](s Stream[T], mapper func(T) Stream[float64]) Stream[float64] {
	return FlatMap(s, mapper)
}

// Reduce folds stream into one value.
// Reduce 将流折叠为单个值。
func (s Stream[T]) Reduce(initial T, operator func(T, T) T) T {
	if operator == nil {
		return initial
	}
	result := initial
	s.seq(func(v T) bool {
		result = operator(result, v)
		return true
	})
	return result
}

// Optional models Java Optional style result.
// Optional 表示 Java Optional 风格结果。
type Optional[T any] struct {
	value   T
	present bool
}

// IsPresent reports whether value exists.
// IsPresent 表示是否存在值。
func (o Optional[T]) IsPresent() bool {
	return o.present
}

// Get returns the optional value and existence flag.
// Get 返回可选值与是否存在。
func (o Optional[T]) Get() (T, bool) {
	return o.value, o.present
}

// OrElse returns default when value not present.
// OrElse 在值不存在时返回默认值。
func (o Optional[T]) OrElse(defaultValue T) T {
	if o.present {
		return o.value
	}
	return defaultValue
}

// ReduceOptional reduces stream without identity.
// ReduceOptional 无初始值折叠流。
func (s Stream[T]) ReduceOptional(operator func(T, T) T) Optional[T] {
	if operator == nil {
		return Optional[T]{}
	}
	first, ok := s.Head()
	if !ok {
		return Optional[T]{}
	}
	result := first
	skippedFirst := false
	s.seq(func(v T) bool {
		if !skippedFirst {
			skippedFirst = true
			return true
		}
		result = operator(result, v)
		return true
	})
	return Optional[T]{value: result, present: true}
}

// ReduceWithCombiner mirrors Java 3-arg reduce signature.
// ReduceWithCombiner 对齐 Java 三参 reduce。
func ReduceWithCombiner[T any, U any](s Stream[T], identity U, accumulator func(U, T) U, combiner func(U, U) U) U {
	if accumulator == nil {
		return identity
	}
	result := identity
	s.seq(func(v T) bool {
		result = accumulator(result, v)
		return true
	})
	if combiner != nil {
		return combiner(identity, result)
	}
	return result
}

// Collect mirrors Java supplier/accumulator/combiner collector.
// Collect 对齐 Java 的 supplier/accumulator/combiner。
func Collect[T any, A any](s Stream[T], supplier func() A, accumulator func(*A, T), combiner func(*A, A)) A {
	if supplier == nil {
		var zero A
		return zero
	}
	result := supplier()
	if accumulator == nil {
		return result
	}
	s.seq(func(v T) bool {
		accumulator(&result, v)
		return true
	})
	if combiner != nil {
		empty := supplier()
		combiner(&result, empty)
	}
	return result
}
