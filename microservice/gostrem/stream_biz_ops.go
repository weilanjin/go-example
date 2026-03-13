package main

// Number defines numeric types for lightweight aggregations.
// Number 定义轻量聚合可用的数值类型。
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// KeyByFn is a Go-style alias of KeyBy.
// KeyByFn 是 KeyBy 的 Go 风格别名。
func KeyByFn[T any, K comparable](s Stream[T], keyFn func(T) K) Stream[Pair[K, T]] {
	return KeyBy(s, keyFn)
}

// ReduceByKeyFn is a Go-style alias of ReduceByKey.
// ReduceByKeyFn 是 ReduceByKey 的 Go 风格别名。
func ReduceByKeyFn[K comparable, V any](s Stream[Pair[K, V]], reducer func(V, V) V) Stream[Pair[K, V]] {
	return ReduceByKey(s, reducer)
}

// SumByKey sums numeric values grouped by key.
// SumByKey 按 key 分组并对数值求和。
func SumByKey[K comparable, N Number](s Stream[Pair[K, N]]) Stream[Pair[K, N]] {
	return ReduceByKey(s, func(a, b N) N { return a + b })
}

// CountByKeyStream counts rows per key and returns stream pairs.
// CountByKeyStream 统计每个 key 的记录数并返回流。
func CountByKeyStream[K comparable, V any](s Stream[Pair[K, V]]) Stream[Pair[K, int64]] {
	m := CountByKey(s)
	out := make([]Pair[K, int64], 0, len(m))
	for k, v := range m {
		out = append(out, Pair[K, int64]{Key: k, Value: v})
	}
	return FromSlice(out)
}

// InnerJoinByKey is a Go-style alias of Join.
// InnerJoinByKey 是 Join 的 Go 风格别名。
func InnerJoinByKey[K comparable, V any, W any](left Stream[Pair[K, V]], right Stream[Pair[K, W]]) Stream[Pair[K, Tuple2[V, W]]] {
	return Join(left, right)
}

// LeftJoinByKey is a Go-style alias of LeftOuterJoin.
// LeftJoinByKey 是 LeftOuterJoin 的 Go 风格别名。
func LeftJoinByKey[K comparable, V any, W any](left Stream[Pair[K, V]], right Stream[Pair[K, W]]) Stream[Pair[K, Tuple2[V, Optional[W]]]] {
	return LeftOuterJoin(left, right)
}

// TumblingWindow groups elements by fixed count windows.
// TumblingWindow 按固定条数进行滚动窗口分组。
func TumblingWindow[T any](s Stream[T], size int) Stream[[]T] {
	return Batch(s, size)
}

// SlidingWindow groups elements by sliding count windows.
// SlidingWindow 按滑动条数进行窗口分组。
func SlidingWindow[T any](s Stream[T], size int, slide int) Stream[[]T] {
	if size <= 0 {
		size = 1
	}
	if slide <= 0 {
		slide = 1
	}
	items := s.Slice()
	if len(items) == 0 {
		return Empty[[]T]()
	}
	out := make([][]T, 0)
	for start := 0; start < len(items); start += slide {
		end := start + size
		if end > len(items) {
			end = len(items)
		}
		if start >= end {
			break
		}
		win := make([]T, end-start)
		copy(win, items[start:end])
		out = append(out, win)
		if end == len(items) {
			break
		}
	}
	return FromSlice(out)
}

// WindowReduce reduces each tumbling count window into one value.
// WindowReduce 将每个滚动条数窗口归约为一个值。
func WindowReduce[T any, R any](s Stream[T], size int, zero R, reduce func(R, T) R) Stream[R] {
	if reduce == nil {
		return Empty[R]()
	}
	wins := TumblingWindow(s, size).Slice()
	out := make([]R, 0, len(wins))
	for i := range wins {
		acc := zero
		for j := range wins[i] {
			acc = reduce(acc, wins[i][j])
		}
		out = append(out, acc)
	}
	return FromSlice(out)
}

// WindowReduceByKey reduces each tumbling window by key.
// WindowReduceByKey 对每个滚动窗口按 key 进行归约。
func WindowReduceByKey[K comparable, V any, R any](s Stream[Pair[K, V]], size int, zero R, reduce func(R, V) R) Stream[Pair[K, R]] {
	if reduce == nil {
		return Empty[Pair[K, R]]()
	}
	wins := TumblingWindow(s, size).Slice()
	out := make([]Pair[K, R], 0)
	for i := range wins {
		acc := make(map[K]R)
		init := make(map[K]bool)
		for j := range wins[i] {
			row := wins[i][j]
			cur := zero
			if init[row.Key] {
				cur = acc[row.Key]
			}
			acc[row.Key] = reduce(cur, row.Value)
			init[row.Key] = true
		}
		for k, v := range acc {
			out = append(out, Pair[K, R]{Key: k, Value: v})
		}
	}
	return FromSlice(out)
}

// Process maps and filters in one pass using ok flag.
// Process 通过 ok 标记在一次遍历中完成映射与过滤。
func Process[T any, R any](s Stream[T], fn func(T) (R, bool)) Stream[R] {
	if fn == nil {
		return Empty[R]()
	}
	seq := func(yield func(R) bool) {
		s.seq(func(v T) bool {
			r, ok := fn(v)
			if !ok {
				return true
			}
			return yield(r)
		})
	}
	return makeStream(seq, -1)
}

// ProcessMany expands one element to many output rows.
// ProcessMany 将一个元素扩展为多个输出元素。
func ProcessMany[T any, R any](s Stream[T], fn func(T) []R) Stream[R] {
	if fn == nil {
		return Empty[R]()
	}
	return FlatMap(s, func(v T) Stream[R] {
		return FromSlice(fn(v))
	})
}

// UpsertByKey materializes latest value per key.
// UpsertByKey 物化每个 key 的最新值。
func UpsertByKey[K comparable, V any](s Stream[Pair[K, V]]) map[K]V {
	return CollectAsMap(s)
}
