package main

import (
	"cmp"
	"math/rand"
	"sort"
)

func Union[T any](left, right Stream[T]) Stream[T] { return Concat(left, right) }

func Intersection[T comparable](left, right Stream[T]) Stream[T] {
	set := make(map[T]struct{})
	right.seq(func(v T) bool { set[v] = struct{}{}; return true })
	return left.Filter(func(v T) bool { _, ok := set[v]; return ok })
}

func Subtract[T comparable](left, right Stream[T]) Stream[T] {
	set := make(map[T]struct{})
	right.seq(func(v T) bool { set[v] = struct{}{}; return true })
	return left.Filter(func(v T) bool { _, ok := set[v]; return !ok })
}

func Cartesian[T any, U any](left Stream[T], right Stream[U]) Stream[Tuple2[T, U]] {
	items := right.Slice()
	seq := func(yield func(Tuple2[T, U]) bool) {
		stop := false
		left.seq(func(v T) bool {
			for i := range items {
				if !yield(Tuple2[T, U]{First: v, Second: items[i]}) {
					stop = true
					return false
				}
			}
			return !stop
		})
	}
	hint := -1
	if left.sizeHint >= 0 && right.sizeHint >= 0 {
		hint = left.sizeHint * right.sizeHint
	}
	return makeStream(seq, hint)
}

func Zip[T any, U any](left Stream[T], right Stream[U]) Stream[Tuple2[T, U]] {
	l := left.Slice()
	r := right.Slice()
	n := len(l)
	if len(r) < n {
		n = len(r)
	}
	seq := func(yield func(Tuple2[T, U]) bool) {
		for i := 0; i < n; i++ {
			if !yield(Tuple2[T, U]{First: l[i], Second: r[i]}) {
				return
			}
		}
	}
	return makeStream(seq, n)
}

func ZipWithIndex[T any](s Stream[T]) Stream[Tuple2[T, int]] {
	seq := func(yield func(Tuple2[T, int]) bool) {
		idx := 0
		s.seq(func(v T) bool {
			ok := yield(Tuple2[T, int]{First: v, Second: idx})
			idx++
			return ok
		})
	}
	return makeStream(seq, s.sizeHint)
}

func Sample[T any](s Stream[T], withReplacement bool, fraction float64, seed int64) Stream[T] {
	if fraction <= 0 {
		return Empty[T]()
	}
	rng := rand.New(rand.NewSource(seed))
	seq := func(yield func(T) bool) {
		s.seq(func(v T) bool {
			if withReplacement {
				for rng.Float64() < fraction {
					if !yield(v) {
						return false
					}
				}
				return true
			}
			if rng.Float64() < fraction {
				return yield(v)
			}
			return true
		})
	}
	return makeStream(seq, -1)
}

func SortBy[T any, K cmp.Ordered](s Stream[T], keyFn func(T) K, ascending bool) Stream[T] {
	if keyFn == nil {
		return s
	}
	if ascending {
		return s.Sorted(func(a, b T) bool { return keyFn(a) < keyFn(b) })
	}
	return s.Sorted(func(a, b T) bool { return keyFn(a) > keyFn(b) })
}

func Glom[T any](s Stream[T], chunkSize int) Stream[[]T] {
	if chunkSize <= 0 {
		chunkSize = 1
	}
	seq := func(yield func([]T) bool) {
		chunk := make([]T, 0, chunkSize)
		flush := func() bool {
			if len(chunk) == 0 {
				return true
			}
			out := make([]T, len(chunk))
			copy(out, chunk)
			chunk = chunk[:0]
			return yield(out)
		}
		s.seq(func(v T) bool {
			chunk = append(chunk, v)
			if len(chunk) < chunkSize {
				return true
			}
			return flush()
		})
		_ = flush()
	}
	return makeStream(seq, -1)
}

func (s Stream[T]) Coalesce(_ int) Stream[T]    { return s }
func (s Stream[T]) Repartition(_ int) Stream[T] { return s }

func MapValues[K comparable, V any, U any](s Stream[Pair[K, V]], mapper func(V) U) Stream[Pair[K, U]] {
	if mapper == nil {
		return Empty[Pair[K, U]]()
	}
	return Map(s, func(p Pair[K, V]) Pair[K, U] { return Pair[K, U]{Key: p.Key, Value: mapper(p.Value)} })
}

func FlatMapValues[K comparable, V any, U any](s Stream[Pair[K, V]], mapper func(V) []U) Stream[Pair[K, U]] {
	if mapper == nil {
		return Empty[Pair[K, U]]()
	}
	return FlatMap(s, func(p Pair[K, V]) Stream[Pair[K, U]] {
		arr := mapper(p.Value)
		mapped := make([]Pair[K, U], 0, len(arr))
		for i := range arr {
			mapped = append(mapped, Pair[K, U]{Key: p.Key, Value: arr[i]})
		}
		return FromSlice(mapped)
	})
}

func Keys[K comparable, V any](s Stream[Pair[K, V]]) Stream[K] {
	return Map(s, func(p Pair[K, V]) K { return p.Key })
}
func Values[K comparable, V any](s Stream[Pair[K, V]]) Stream[V] {
	return Map(s, func(p Pair[K, V]) V { return p.Value })
}

func GroupBy[T any, K comparable](s Stream[T], classifier func(T) K) Stream[Pair[K, []T]] {
	grouped := make(map[K][]T)
	s.seq(func(v T) bool {
		k := classifier(v)
		grouped[k] = append(grouped[k], v)
		return true
	})
	res := make([]Pair[K, []T], 0, len(grouped))
	for k, v := range grouped {
		res = append(res, Pair[K, []T]{Key: k, Value: v})
	}
	return FromSlice(res)
}

func GroupByKey[K comparable, V any](s Stream[Pair[K, V]]) Stream[Pair[K, []V]] {
	grouped := GroupBy(s, func(p Pair[K, V]) K { return p.Key })
	return Map(grouped, func(p Pair[K, []Pair[K, V]]) Pair[K, []V] {
		vals := make([]V, 0, len(p.Value))
		for i := range p.Value {
			vals = append(vals, p.Value[i].Value)
		}
		return Pair[K, []V]{Key: p.Key, Value: vals}
	})
}

func ReduceByKey[K comparable, V any](s Stream[Pair[K, V]], reducer func(V, V) V) Stream[Pair[K, V]] {
	if reducer == nil {
		return s
	}
	acc := make(map[K]V)
	init := make(map[K]bool)
	s.seq(func(v Pair[K, V]) bool {
		if !init[v.Key] {
			acc[v.Key] = v.Value
			init[v.Key] = true
			return true
		}
		acc[v.Key] = reducer(acc[v.Key], v.Value)
		return true
	})
	out := make([]Pair[K, V], 0, len(acc))
	for k, v := range acc {
		out = append(out, Pair[K, V]{Key: k, Value: v})
	}
	return FromSlice(out)
}

func FoldByKey[K comparable, V any](s Stream[Pair[K, V]], zero V, op func(V, V) V) Stream[Pair[K, V]] {
	if op == nil {
		return s
	}
	acc := make(map[K]V)
	s.seq(func(v Pair[K, V]) bool {
		cur, ok := acc[v.Key]
		if !ok {
			cur = zero
		}
		acc[v.Key] = op(cur, v.Value)
		return true
	})
	out := make([]Pair[K, V], 0, len(acc))
	for k, v := range acc {
		out = append(out, Pair[K, V]{Key: k, Value: v})
	}
	return FromSlice(out)
}

func SortByKey[K cmp.Ordered, V any](s Stream[Pair[K, V]], ascending bool) Stream[Pair[K, V]] {
	if ascending {
		return s.Sorted(func(a, b Pair[K, V]) bool { return a.Key < b.Key })
	}
	return s.Sorted(func(a, b Pair[K, V]) bool { return a.Key > b.Key })
}

func CountByValue[T comparable](s Stream[T]) map[T]int64 {
	out := make(map[T]int64)
	s.seq(func(v T) bool { out[v]++; return true })
	return out
}

func CountByKey[K comparable, V any](s Stream[Pair[K, V]]) map[K]int64 {
	out := make(map[K]int64)
	s.seq(func(v Pair[K, V]) bool { out[v.Key]++; return true })
	return out
}

func (s Stream[T]) Take(n int) []T   { return s.Limit(n).Slice() }
func (s Stream[T]) First() (T, bool) { return s.Head() }

func TakeOrdered[T cmp.Ordered](s Stream[T], n int) []T {
	if n <= 0 {
		return []T{}
	}
	items := s.Slice()
	sort.Slice(items, func(i, j int) bool { return items[i] < items[j] })
	if n > len(items) {
		n = len(items)
	}
	return items[:n]
}

func Top[T cmp.Ordered](s Stream[T], n int) []T {
	if n <= 0 {
		return []T{}
	}
	items := s.Slice()
	sort.Slice(items, func(i, j int) bool { return items[i] > items[j] })
	if n > len(items) {
		n = len(items)
	}
	return items[:n]
}

func KeyBy[T any, K comparable](s Stream[T], keyFn func(T) K) Stream[Pair[K, T]] {
	if keyFn == nil {
		return Empty[Pair[K, T]]()
	}
	return Map(s, func(v T) Pair[K, T] { return Pair[K, T]{Key: keyFn(v), Value: v} })
}

func Lookup[K comparable, V any](s Stream[Pair[K, V]], key K) []V {
	out := make([]V, 0)
	s.seq(func(p Pair[K, V]) bool {
		if p.Key == key {
			out = append(out, p.Value)
		}
		return true
	})
	return out
}

func CollectAsMap[K comparable, V any](s Stream[Pair[K, V]]) map[K]V {
	out := make(map[K]V)
	s.seq(func(p Pair[K, V]) bool { out[p.Key] = p.Value; return true })
	return out
}

func Join[K comparable, V any, W any](left Stream[Pair[K, V]], right Stream[Pair[K, W]]) Stream[Pair[K, Tuple2[V, W]]] {
	rightMap := make(map[K][]W)
	right.seq(func(p Pair[K, W]) bool { rightMap[p.Key] = append(rightMap[p.Key], p.Value); return true })
	seq := func(yield func(Pair[K, Tuple2[V, W]]) bool) {
		left.seq(func(lp Pair[K, V]) bool {
			vals, ok := rightMap[lp.Key]
			if !ok {
				return true
			}
			for i := range vals {
				if !yield(Pair[K, Tuple2[V, W]]{Key: lp.Key, Value: Tuple2[V, W]{First: lp.Value, Second: vals[i]}}) {
					return false
				}
			}
			return true
		})
	}
	return makeStream(seq, -1)
}

func LeftOuterJoin[K comparable, V any, W any](left Stream[Pair[K, V]], right Stream[Pair[K, W]]) Stream[Pair[K, Tuple2[V, Optional[W]]]] {
	rightMap := make(map[K][]W)
	right.seq(func(p Pair[K, W]) bool { rightMap[p.Key] = append(rightMap[p.Key], p.Value); return true })
	seq := func(yield func(Pair[K, Tuple2[V, Optional[W]]]) bool) {
		left.seq(func(lp Pair[K, V]) bool {
			vals, ok := rightMap[lp.Key]
			if !ok || len(vals) == 0 {
				return yield(Pair[K, Tuple2[V, Optional[W]]]{Key: lp.Key, Value: Tuple2[V, Optional[W]]{First: lp.Value, Second: Optional[W]{}}})
			}
			for i := range vals {
				if !yield(Pair[K, Tuple2[V, Optional[W]]]{Key: lp.Key, Value: Tuple2[V, Optional[W]]{First: lp.Value, Second: Optional[W]{value: vals[i], present: true}}}) {
					return false
				}
			}
			return true
		})
	}
	return makeStream(seq, -1)
}

func RightOuterJoin[K comparable, V any, W any](left Stream[Pair[K, V]], right Stream[Pair[K, W]]) Stream[Pair[K, Tuple2[Optional[V], W]]] {
	leftMap := make(map[K][]V)
	left.seq(func(p Pair[K, V]) bool { leftMap[p.Key] = append(leftMap[p.Key], p.Value); return true })
	seq := func(yield func(Pair[K, Tuple2[Optional[V], W]]) bool) {
		right.seq(func(rp Pair[K, W]) bool {
			vals, ok := leftMap[rp.Key]
			if !ok || len(vals) == 0 {
				return yield(Pair[K, Tuple2[Optional[V], W]]{Key: rp.Key, Value: Tuple2[Optional[V], W]{First: Optional[V]{}, Second: rp.Value}})
			}
			for i := range vals {
				if !yield(Pair[K, Tuple2[Optional[V], W]]{Key: rp.Key, Value: Tuple2[Optional[V], W]{First: Optional[V]{value: vals[i], present: true}, Second: rp.Value}}) {
					return false
				}
			}
			return true
		})
	}
	return makeStream(seq, -1)
}

func FullOuterJoin[K comparable, V any, W any](left Stream[Pair[K, V]], right Stream[Pair[K, W]]) Stream[Pair[K, Tuple2[Optional[V], Optional[W]]]] {
	leftMap := make(map[K][]V)
	left.seq(func(p Pair[K, V]) bool { leftMap[p.Key] = append(leftMap[p.Key], p.Value); return true })
	rightMap := make(map[K][]W)
	right.seq(func(p Pair[K, W]) bool { rightMap[p.Key] = append(rightMap[p.Key], p.Value); return true })
	seq := func(yield func(Pair[K, Tuple2[Optional[V], Optional[W]]]) bool) {
		keys := make(map[K]struct{})
		for k := range leftMap {
			keys[k] = struct{}{}
		}
		for k := range rightMap {
			keys[k] = struct{}{}
		}
		for k := range keys {
			lvals := leftMap[k]
			rvals := rightMap[k]
			if len(lvals) == 0 {
				for i := range rvals {
					if !yield(Pair[K, Tuple2[Optional[V], Optional[W]]]{Key: k, Value: Tuple2[Optional[V], Optional[W]]{First: Optional[V]{}, Second: Optional[W]{value: rvals[i], present: true}}}) {
						return
					}
				}
				continue
			}
			if len(rvals) == 0 {
				for i := range lvals {
					if !yield(Pair[K, Tuple2[Optional[V], Optional[W]]]{Key: k, Value: Tuple2[Optional[V], Optional[W]]{First: Optional[V]{value: lvals[i], present: true}, Second: Optional[W]{}}}) {
						return
					}
				}
				continue
			}
			for i := range lvals {
				for j := range rvals {
					if !yield(Pair[K, Tuple2[Optional[V], Optional[W]]]{Key: k, Value: Tuple2[Optional[V], Optional[W]]{First: Optional[V]{value: lvals[i], present: true}, Second: Optional[W]{value: rvals[j], present: true}}}) {
						return
					}
				}
			}
		}
	}
	return makeStream(seq, -1)
}

func Cogroup[K comparable, V any, W any](left Stream[Pair[K, V]], right Stream[Pair[K, W]]) Stream[Pair[K, Tuple2[[]V, []W]]] {
	leftMap := make(map[K][]V)
	left.seq(func(p Pair[K, V]) bool { leftMap[p.Key] = append(leftMap[p.Key], p.Value); return true })
	rightMap := make(map[K][]W)
	right.seq(func(p Pair[K, W]) bool { rightMap[p.Key] = append(rightMap[p.Key], p.Value); return true })
	out := make([]Pair[K, Tuple2[[]V, []W]], 0)
	keys := make(map[K]struct{})
	for k := range leftMap {
		keys[k] = struct{}{}
	}
	for k := range rightMap {
		keys[k] = struct{}{}
	}
	for k := range keys {
		out = append(out, Pair[K, Tuple2[[]V, []W]]{Key: k, Value: Tuple2[[]V, []W]{First: leftMap[k], Second: rightMap[k]}})
	}
	return FromSlice(out)
}

func AggregateByKey[K comparable, V any, U any](s Stream[Pair[K, V]], zero U, seqOp func(U, V) U, combOp func(U, U) U) Stream[Pair[K, U]] {
	if seqOp == nil {
		return Empty[Pair[K, U]]()
	}
	acc := make(map[K]U)
	s.seq(func(p Pair[K, V]) bool {
		cur, ok := acc[p.Key]
		if !ok {
			cur = zero
		}
		acc[p.Key] = seqOp(cur, p.Value)
		return true
	})
	_ = combOp
	out := make([]Pair[K, U], 0, len(acc))
	for k, v := range acc {
		out = append(out, Pair[K, U]{Key: k, Value: v})
	}
	return FromSlice(out)
}

func CombineByKey[K comparable, V any, C any](s Stream[Pair[K, V]], createCombiner func(V) C, mergeValue func(C, V) C, mergeCombiners func(C, C) C) Stream[Pair[K, C]] {
	if createCombiner == nil || mergeValue == nil {
		return Empty[Pair[K, C]]()
	}
	acc := make(map[K]C)
	init := make(map[K]bool)
	s.seq(func(p Pair[K, V]) bool {
		if !init[p.Key] {
			acc[p.Key] = createCombiner(p.Value)
			init[p.Key] = true
			return true
		}
		acc[p.Key] = mergeValue(acc[p.Key], p.Value)
		return true
	})
	_ = mergeCombiners
	out := make([]Pair[K, C], 0, len(acc))
	for k, v := range acc {
		out = append(out, Pair[K, C]{Key: k, Value: v})
	}
	return FromSlice(out)
}

func MapPartitions[T any, R any](s Stream[T], fn func([]T) []R) Stream[R] {
	if fn == nil {
		return Empty[R]()
	}
	return FromSlice(fn(s.Slice()))
}

func MapPartitionsWithIndex[T any, R any](s Stream[T], fn func(int, []T) []R) Stream[R] {
	if fn == nil {
		return Empty[R]()
	}
	return FromSlice(fn(0, s.Slice()))
}

func ForEachPartition[T any](s Stream[T], consumer func([]T)) {
	if consumer == nil {
		return
	}
	consumer(s.Slice())
}
