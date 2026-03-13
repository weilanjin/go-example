package main

import "strings"

// Collector models Java Collectors contract in Go generic style.
// Collector 用 Go 泛型表达 Java Collectors 契约。
type Collector[T any, A any, R any] struct {
	Supplier    func() A
	Accumulator func(*A, T)
	Finisher    func(A) R
}

// CollectWith collects stream by collector contract.
// CollectWith 使用 collector 契约收集流。
func CollectWith[T any, A any, R any](s Stream[T], c Collector[T, A, R]) R {
	if c.Supplier == nil {
		var zero R
		return zero
	}
	acc := c.Supplier()
	if c.Accumulator != nil {
		s.seq(func(v T) bool {
			c.Accumulator(&acc, v)
			return true
		})
	}
	if c.Finisher != nil {
		return c.Finisher(acc)
	}
	var zero R
	return zero
}

// ToSliceCollector returns collector that gathers elements into slice.
// ToSliceCollector 返回收集为切片的 collector。
func ToSliceCollector[T any]() Collector[T, []T, []T] {
	return Collector[T, []T, []T]{
		Supplier: func() []T { return make([]T, 0) },
		Accumulator: func(dst *[]T, v T) {
			*dst = append(*dst, v)
		},
		Finisher: func(acc []T) []T { return acc },
	}
}

// ToSetCollector returns collector that gathers elements into set map.
// ToSetCollector 返回收集为集合 map 的 collector。
func ToSetCollector[T comparable]() Collector[T, map[T]struct{}, map[T]struct{}] {
	return Collector[T, map[T]struct{}, map[T]struct{}]{
		Supplier: func() map[T]struct{} { return map[T]struct{}{} },
		Accumulator: func(dst *map[T]struct{}, v T) {
			(*dst)[v] = struct{}{}
		},
		Finisher: func(acc map[T]struct{}) map[T]struct{} { return acc },
	}
}

// CountingCollector returns collector that counts elements.
// CountingCollector 返回计数 collector。
func CountingCollector[T any]() Collector[T, int64, int64] {
	return Collector[T, int64, int64]{
		Supplier: func() int64 { return 0 },
		Accumulator: func(dst *int64, _ T) {
			*dst = *dst + 1
		},
		Finisher: func(acc int64) int64 { return acc },
	}
}

// JoiningCollector joins string elements.
// JoiningCollector 拼接字符串元素。
func JoiningCollector(delimiter, prefix, suffix string) Collector[string, []string, string] {
	return Collector[string, []string, string]{
		Supplier: func() []string { return make([]string, 0) },
		Accumulator: func(dst *[]string, v string) {
			*dst = append(*dst, v)
		},
		Finisher: func(acc []string) string {
			return prefix + strings.Join(acc, delimiter) + suffix
		},
	}
}

// MappingCollector adapts input element type for downstream collector.
// MappingCollector 为下游 collector 适配输入类型。
func MappingCollector[T any, U any, A any, R any](mapper func(T) U, downstream Collector[U, A, R]) Collector[T, A, R] {
	if mapper == nil || downstream.Supplier == nil {
		return Collector[T, A, R]{}
	}
	return Collector[T, A, R]{
		Supplier: downstream.Supplier,
		Accumulator: func(dst *A, v T) {
			if downstream.Accumulator == nil {
				return
			}
			downstream.Accumulator(dst, mapper(v))
		},
		Finisher: downstream.Finisher,
	}
}

// FilteringCollector forwards only matched elements to downstream collector.
// FilteringCollector 仅将满足条件的元素发送到下游 collector。
func FilteringCollector[T any, A any, R any](predicate func(T) bool, downstream Collector[T, A, R]) Collector[T, A, R] {
	if predicate == nil || downstream.Supplier == nil {
		return Collector[T, A, R]{}
	}
	return Collector[T, A, R]{
		Supplier: downstream.Supplier,
		Accumulator: func(dst *A, v T) {
			if !predicate(v) || downstream.Accumulator == nil {
				return
			}
			downstream.Accumulator(dst, v)
		},
		Finisher: downstream.Finisher,
	}
}

// FlatMappingCollector expands each input into many downstream values.
// FlatMappingCollector 将输入展开为多个下游值。
func FlatMappingCollector[T any, U any, A any, R any](mapper func(T) []U, downstream Collector[U, A, R]) Collector[T, A, R] {
	if mapper == nil || downstream.Supplier == nil {
		return Collector[T, A, R]{}
	}
	return Collector[T, A, R]{
		Supplier: downstream.Supplier,
		Accumulator: func(dst *A, v T) {
			if downstream.Accumulator == nil {
				return
			}
			items := mapper(v)
			for i := range items {
				downstream.Accumulator(dst, items[i])
			}
		},
		Finisher: downstream.Finisher,
	}
}

// GroupingByCollector groups values by key into map[K][]T.
// GroupingByCollector 按 key 分组为 map[K][]T。
func GroupingByCollector[T any, K comparable](classifier func(T) K) Collector[T, map[K][]T, map[K][]T] {
	if classifier == nil {
		return Collector[T, map[K][]T, map[K][]T]{}
	}
	return Collector[T, map[K][]T, map[K][]T]{
		Supplier: func() map[K][]T { return map[K][]T{} },
		Accumulator: func(dst *map[K][]T, v T) {
			k := classifier(v)
			(*dst)[k] = append((*dst)[k], v)
		},
		Finisher: func(acc map[K][]T) map[K][]T { return acc },
	}
}

// GroupingByMappingCollector groups by key and maps each value before appending.
// GroupingByMappingCollector 按 key 分组并映射每个值。
func GroupingByMappingCollector[T any, K comparable, V any](classifier func(T) K, mapper func(T) V) Collector[T, map[K][]V, map[K][]V] {
	if classifier == nil || mapper == nil {
		return Collector[T, map[K][]V, map[K][]V]{}
	}
	return Collector[T, map[K][]V, map[K][]V]{
		Supplier: func() map[K][]V { return map[K][]V{} },
		Accumulator: func(dst *map[K][]V, v T) {
			k := classifier(v)
			(*dst)[k] = append((*dst)[k], mapper(v))
		},
		Finisher: func(acc map[K][]V) map[K][]V { return acc },
	}
}

// PartitioningByCollector partitions values into true/false buckets.
// PartitioningByCollector 将值划分为 true/false 两个桶。
func PartitioningByCollector[T any](predicate func(T) bool) Collector[T, map[bool][]T, map[bool][]T] {
	if predicate == nil {
		return Collector[T, map[bool][]T, map[bool][]T]{}
	}
	return Collector[T, map[bool][]T, map[bool][]T]{
		Supplier: func() map[bool][]T {
			return map[bool][]T{true: make([]T, 0), false: make([]T, 0)}
		},
		Accumulator: func(dst *map[bool][]T, v T) {
			k := predicate(v)
			(*dst)[k] = append((*dst)[k], v)
		},
		Finisher: func(acc map[bool][]T) map[bool][]T { return acc },
	}
}

// ToMapCollector collects values into map with merge strategy on duplicate keys.
// ToMapCollector 收集到 map，并支持重复 key 合并策略。
func ToMapCollector[T any, K comparable, V any](keyMapper func(T) K, valueMapper func(T) V, mergeFn func(existing, incoming V) V) Collector[T, map[K]V, map[K]V] {
	if keyMapper == nil || valueMapper == nil {
		return Collector[T, map[K]V, map[K]V]{}
	}
	return Collector[T, map[K]V, map[K]V]{
		Supplier: func() map[K]V { return map[K]V{} },
		Accumulator: func(dst *map[K]V, v T) {
			k := keyMapper(v)
			nv := valueMapper(v)
			if ov, ok := (*dst)[k]; ok && mergeFn != nil {
				(*dst)[k] = mergeFn(ov, nv)
				return
			}
			(*dst)[k] = nv
		},
		Finisher: func(acc map[K]V) map[K]V { return acc },
	}
}

// MinByCollector collects minimal element by comparator.
// MinByCollector 按比较器收集最小元素。
func MinByCollector[T any](less func(a, b T) bool) Collector[T, Optional[T], Optional[T]] {
	if less == nil {
		return Collector[T, Optional[T], Optional[T]]{}
	}
	return Collector[T, Optional[T], Optional[T]]{
		Supplier: func() Optional[T] { return Optional[T]{} },
		Accumulator: func(dst *Optional[T], v T) {
			if !dst.present || less(v, dst.value) {
				dst.value = v
				dst.present = true
			}
		},
		Finisher: func(acc Optional[T]) Optional[T] { return acc },
	}
}

// MaxByCollector collects maximal element by comparator.
// MaxByCollector 按比较器收集最大元素。
func MaxByCollector[T any](less func(a, b T) bool) Collector[T, Optional[T], Optional[T]] {
	if less == nil {
		return Collector[T, Optional[T], Optional[T]]{}
	}
	return Collector[T, Optional[T], Optional[T]]{
		Supplier: func() Optional[T] { return Optional[T]{} },
		Accumulator: func(dst *Optional[T], v T) {
			if !dst.present || less(dst.value, v) {
				dst.value = v
				dst.present = true
			}
		},
		Finisher: func(acc Optional[T]) Optional[T] { return acc },
	}
}

// GroupingByDownstreamCollector groups by key and applies downstream collector per bucket.
// GroupingByDownstreamCollector 按 key 分组并对每组应用下游 collector。
func GroupingByDownstreamCollector[T any, K comparable, A any, R any](classifier func(T) K, downstream Collector[T, A, R]) Collector[T, map[K]A, map[K]R] {
	if classifier == nil || downstream.Supplier == nil {
		return Collector[T, map[K]A, map[K]R]{}
	}
	return Collector[T, map[K]A, map[K]R]{
		Supplier: func() map[K]A { return map[K]A{} },
		Accumulator: func(dst *map[K]A, v T) {
			k := classifier(v)
			bucket, ok := (*dst)[k]
			if !ok {
				bucket = downstream.Supplier()
			}
			if downstream.Accumulator != nil {
				downstream.Accumulator(&bucket, v)
			}
			(*dst)[k] = bucket
		},
		Finisher: func(acc map[K]A) map[K]R {
			out := make(map[K]R, len(acc))
			for k, v := range acc {
				if downstream.Finisher != nil {
					out[k] = downstream.Finisher(v)
					continue
				}
				var zero R
				out[k] = zero
			}
			return out
		},
	}
}

type teeAccumulator[A1 any, A2 any] struct {
	left  A1
	right A2
}

// TeeingCollector routes stream to two collectors and merges final results.
// TeeingCollector 将流送入两个 collector 并合并最终结果。
func TeeingCollector[T any, A1 any, R1 any, A2 any, R2 any, R any](
	left Collector[T, A1, R1],
	right Collector[T, A2, R2],
	merger func(R1, R2) R,
) Collector[T, teeAccumulator[A1, A2], R] {
	if left.Supplier == nil || right.Supplier == nil || merger == nil {
		return Collector[T, teeAccumulator[A1, A2], R]{}
	}
	return Collector[T, teeAccumulator[A1, A2], R]{
		Supplier: func() teeAccumulator[A1, A2] {
			return teeAccumulator[A1, A2]{
				left:  left.Supplier(),
				right: right.Supplier(),
			}
		},
		Accumulator: func(dst *teeAccumulator[A1, A2], v T) {
			if left.Accumulator != nil {
				left.Accumulator(&dst.left, v)
			}
			if right.Accumulator != nil {
				right.Accumulator(&dst.right, v)
			}
		},
		Finisher: func(acc teeAccumulator[A1, A2]) R {
			var lres R1
			var rres R2
			if left.Finisher != nil {
				lres = left.Finisher(acc.left)
			}
			if right.Finisher != nil {
				rres = right.Finisher(acc.right)
			}
			return merger(lres, rres)
		},
	}
}

// CollectingAndThen wraps downstream finisher with final transform.
// CollectingAndThen 使用最终转换包装下游 finisher。
func CollectingAndThen[T any, A any, R any, RR any](downstream Collector[T, A, R], finisher func(R) RR) Collector[T, A, RR] {
	if downstream.Supplier == nil || finisher == nil {
		return Collector[T, A, RR]{}
	}
	return Collector[T, A, RR]{
		Supplier:    downstream.Supplier,
		Accumulator: downstream.Accumulator,
		Finisher: func(acc A) RR {
			var mid R
			if downstream.Finisher != nil {
				mid = downstream.Finisher(acc)
			}
			return finisher(mid)
		},
	}
}

// ReducingCollector reduces mapped values by operator with identity.
// ReducingCollector 使用 identity 与操作符进行归约。
func ReducingCollector[T any, U any](identity U, mapper func(T) U, op func(U, U) U) Collector[T, U, U] {
	if mapper == nil || op == nil {
		return Collector[T, U, U]{}
	}
	return Collector[T, U, U]{
		Supplier: func() U { return identity },
		Accumulator: func(dst *U, v T) {
			*dst = op(*dst, mapper(v))
		},
		Finisher: func(acc U) U { return acc },
	}
}

// SummingIntCollector sums int-mapped values into int64.
// SummingIntCollector 对 int 映射值求和，结果为 int64。
func SummingIntCollector[T any](mapper func(T) int) Collector[T, int64, int64] {
	if mapper == nil {
		return Collector[T, int64, int64]{}
	}
	return Collector[T, int64, int64]{
		Supplier: func() int64 { return 0 },
		Accumulator: func(dst *int64, v T) {
			*dst += int64(mapper(v))
		},
		Finisher: func(acc int64) int64 { return acc },
	}
}

// SummingInt64Collector sums int64-mapped values.
// SummingInt64Collector 对 int64 映射值求和。
func SummingInt64Collector[T any](mapper func(T) int64) Collector[T, int64, int64] {
	if mapper == nil {
		return Collector[T, int64, int64]{}
	}
	return Collector[T, int64, int64]{
		Supplier: func() int64 { return 0 },
		Accumulator: func(dst *int64, v T) {
			*dst += mapper(v)
		},
		Finisher: func(acc int64) int64 { return acc },
	}
}

// SummingFloat64Collector sums float64-mapped values.
// SummingFloat64Collector 对 float64 映射值求和。
func SummingFloat64Collector[T any](mapper func(T) float64) Collector[T, float64, float64] {
	if mapper == nil {
		return Collector[T, float64, float64]{}
	}
	return Collector[T, float64, float64]{
		Supplier: func() float64 { return 0 },
		Accumulator: func(dst *float64, v T) {
			*dst += mapper(v)
		},
		Finisher: func(acc float64) float64 { return acc },
	}
}

type averagingState struct {
	sum   float64
	count int64
}

// AveragingIntCollector computes average of int-mapped values.
// AveragingIntCollector 计算 int 映射值的平均值。
func AveragingIntCollector[T any](mapper func(T) int) Collector[T, averagingState, float64] {
	if mapper == nil {
		return Collector[T, averagingState, float64]{}
	}
	return Collector[T, averagingState, float64]{
		Supplier: func() averagingState { return averagingState{} },
		Accumulator: func(dst *averagingState, v T) {
			dst.sum += float64(mapper(v))
			dst.count++
		},
		Finisher: func(acc averagingState) float64 {
			if acc.count == 0 {
				return 0
			}
			return acc.sum / float64(acc.count)
		},
	}
}

// AveragingInt64Collector computes average of int64-mapped values.
// AveragingInt64Collector 计算 int64 映射值的平均值。
func AveragingInt64Collector[T any](mapper func(T) int64) Collector[T, averagingState, float64] {
	if mapper == nil {
		return Collector[T, averagingState, float64]{}
	}
	return Collector[T, averagingState, float64]{
		Supplier: func() averagingState { return averagingState{} },
		Accumulator: func(dst *averagingState, v T) {
			dst.sum += float64(mapper(v))
			dst.count++
		},
		Finisher: func(acc averagingState) float64 {
			if acc.count == 0 {
				return 0
			}
			return acc.sum / float64(acc.count)
		},
	}
}

// AveragingFloat64Collector computes average of float64-mapped values.
// AveragingFloat64Collector 计算 float64 映射值的平均值。
func AveragingFloat64Collector[T any](mapper func(T) float64) Collector[T, averagingState, float64] {
	if mapper == nil {
		return Collector[T, averagingState, float64]{}
	}
	return Collector[T, averagingState, float64]{
		Supplier: func() averagingState { return averagingState{} },
		Accumulator: func(dst *averagingState, v T) {
			dst.sum += mapper(v)
			dst.count++
		},
		Finisher: func(acc averagingState) float64 {
			if acc.count == 0 {
				return 0
			}
			return acc.sum / float64(acc.count)
		},
	}
}

// IntSummaryStatistics is summary for int mapped values.
// IntSummaryStatistics 是 int 映射值的统计结果。
type IntSummaryStatistics struct {
	Count int64
	Sum   int64
	Min   int
	Max   int
	init  bool
}

// Average returns average value.
// Average 返回平均值。
func (s IntSummaryStatistics) Average() float64 {
	if s.Count == 0 {
		return 0
	}
	return float64(s.Sum) / float64(s.Count)
}

// LongSummaryStatistics is summary for int64 mapped values.
// LongSummaryStatistics 是 int64 映射值的统计结果。
type LongSummaryStatistics struct {
	Count int64
	Sum   int64
	Min   int64
	Max   int64
	init  bool
}

// Average returns average value.
// Average 返回平均值。
func (s LongSummaryStatistics) Average() float64 {
	if s.Count == 0 {
		return 0
	}
	return float64(s.Sum) / float64(s.Count)
}

// DoubleSummaryStatistics is summary for float64 mapped values.
// DoubleSummaryStatistics 是 float64 映射值的统计结果。
type DoubleSummaryStatistics struct {
	Count int64
	Sum   float64
	Min   float64
	Max   float64
	init  bool
}

// Average returns average value.
// Average 返回平均值。
func (s DoubleSummaryStatistics) Average() float64 {
	if s.Count == 0 {
		return 0
	}
	return s.Sum / float64(s.Count)
}

// SummarizingIntCollector summarizes int-mapped values.
// SummarizingIntCollector 汇总 int 映射值统计信息。
func SummarizingIntCollector[T any](mapper func(T) int) Collector[T, IntSummaryStatistics, IntSummaryStatistics] {
	if mapper == nil {
		return Collector[T, IntSummaryStatistics, IntSummaryStatistics]{}
	}
	return Collector[T, IntSummaryStatistics, IntSummaryStatistics]{
		Supplier: func() IntSummaryStatistics { return IntSummaryStatistics{} },
		Accumulator: func(dst *IntSummaryStatistics, v T) {
			n := mapper(v)
			dst.Count++
			dst.Sum += int64(n)
			if !dst.init {
				dst.Min, dst.Max, dst.init = n, n, true
				return
			}
			if n < dst.Min {
				dst.Min = n
			}
			if n > dst.Max {
				dst.Max = n
			}
		},
		Finisher: func(acc IntSummaryStatistics) IntSummaryStatistics { return acc },
	}
}

// SummarizingInt64Collector summarizes int64-mapped values.
// SummarizingInt64Collector 汇总 int64 映射值统计信息。
func SummarizingInt64Collector[T any](mapper func(T) int64) Collector[T, LongSummaryStatistics, LongSummaryStatistics] {
	if mapper == nil {
		return Collector[T, LongSummaryStatistics, LongSummaryStatistics]{}
	}
	return Collector[T, LongSummaryStatistics, LongSummaryStatistics]{
		Supplier: func() LongSummaryStatistics { return LongSummaryStatistics{} },
		Accumulator: func(dst *LongSummaryStatistics, v T) {
			n := mapper(v)
			dst.Count++
			dst.Sum += n
			if !dst.init {
				dst.Min, dst.Max, dst.init = n, n, true
				return
			}
			if n < dst.Min {
				dst.Min = n
			}
			if n > dst.Max {
				dst.Max = n
			}
		},
		Finisher: func(acc LongSummaryStatistics) LongSummaryStatistics { return acc },
	}
}

// SummarizingFloat64Collector summarizes float64-mapped values.
// SummarizingFloat64Collector 汇总 float64 映射值统计信息。
func SummarizingFloat64Collector[T any](mapper func(T) float64) Collector[T, DoubleSummaryStatistics, DoubleSummaryStatistics] {
	if mapper == nil {
		return Collector[T, DoubleSummaryStatistics, DoubleSummaryStatistics]{}
	}
	return Collector[T, DoubleSummaryStatistics, DoubleSummaryStatistics]{
		Supplier: func() DoubleSummaryStatistics { return DoubleSummaryStatistics{} },
		Accumulator: func(dst *DoubleSummaryStatistics, v T) {
			n := mapper(v)
			dst.Count++
			dst.Sum += n
			if !dst.init {
				dst.Min, dst.Max, dst.init = n, n, true
				return
			}
			if n < dst.Min {
				dst.Min = n
			}
			if n > dst.Max {
				dst.Max = n
			}
		},
		Finisher: func(acc DoubleSummaryStatistics) DoubleSummaryStatistics { return acc },
	}
}
