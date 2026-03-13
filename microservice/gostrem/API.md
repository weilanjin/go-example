# gostrem API 文档

本文档覆盖 microservice/gostrem 包当前导出的主要 API，按模块组织：Core、Terminal、Collectors、Spark-like。

Go 风格优先入口（推荐）：
- 构造：From / FromPointer / FromSlice / Chain
- 过滤与转换：Where / Filter / MapFn / FlatMapFn / Unique / UniqueBy
- 终止：Slice / All / Len / Head / First / Any / Every / None / Each
- 并发：MapPar / MapParUnordered / FlatMapPar / EachPar
- context 与错误：FromCtx / WithCtx / TryMap / TryFlatMap / TryCollect
- 窗口：Window / Batch

## 1. Core API

### 1.1 核心类型

- Stream[T any]
  - 惰性流抽象，内部基于 iter.Seq[T]。
- Tuple2[A any, B any]
  - Spark 风格二元组。
- Pair[K comparable, V any]
  - Key-Value 结构。
- Builder[T any]
  - 构建流的可变构建器。
- Optional[T any]
  - Java Optional 风格的结果封装。

### 1.2 工厂函数

- Empty[T any]() Stream[T]
- Of[T any](values ...T) Stream[T]
- OfNullable[T any](value *T) Stream[T]
- FromSlice[T any](values []T) Stream[T]
- Concat[T any](left, right Stream[T]) Stream[T]
- Generate[T any](supplier func() T) Stream[T]
- Iterate[T any](seed T, next func(T) T) Stream[T]
- IterateWhile[T any](seed T, hasNext func(T) bool, next func(T) T) Stream[T]
- NewBuilder[T any]() *Builder[T]

### 1.3 Builder 方法

- Add(value T) *Builder[T]
- Build() Stream[T]

### 1.4 中间操作

- (s Stream[T]) Filter(predicate func(T) bool) Stream[T]
- (s Stream[T]) Peek(action func(T)) Stream[T]
- (s Stream[T]) Limit(n int) Stream[T]
- (s Stream[T]) Skip(n int) Stream[T]
- (s Stream[T]) TakeWhile(predicate func(T) bool) Stream[T]
- (s Stream[T]) DropWhile(predicate func(T) bool) Stream[T]
- Distinct[T comparable](s Stream[T]) Stream[T]
- DistinctBy[T any, K comparable](s Stream[T], keyFn func(T) K) Stream[T]
- (s Stream[T]) Sorted(less func(a, b T) bool) Stream[T]
- SortedByKey[T any, K cmp.Ordered](s Stream[T], keyFn func(T) K) Stream[T]
- Map[T any, R any](s Stream[T], mapper func(T) R) Stream[R]
- FlatMap[T any, R any](s Stream[T], mapper func(T) Stream[R]) Stream[R]
- MapToInt[T any](s Stream[T], mapper func(T) int) Stream[int]
- MapToInt64[T any](s Stream[T], mapper func(T) int64) Stream[int64]
- MapToFloat64[T any](s Stream[T], mapper func(T) float64) Stream[float64]
- FlatMapToInt[T any](s Stream[T], mapper func(T) Stream[int]) Stream[int]
- FlatMapToInt64[T any](s Stream[T], mapper func(T) Stream[int64]) Stream[int64]
- FlatMapToFloat64[T any](s Stream[T], mapper func(T) Stream[float64]) Stream[float64]

### 1.5 归约与收集（Core）

- (s Stream[T]) Reduce(initial T, operator func(T, T) T) T
- (s Stream[T]) ReduceOptional(operator func(T, T) T) Optional[T]
- ReduceWithCombiner[T any, U any](s Stream[T], identity U, accumulator func(U, T) U, combiner func(U, U) U) U
- Collect[T any, A any](s Stream[T], supplier func() A, accumulator func(*A, T), combiner func(*A, A)) A

### 1.6 Optional 方法

- (o Optional[T]) IsPresent() bool
- (o Optional[T]) Get() (T, bool)
- (o Optional[T]) OrElse(defaultValue T) T

## 2. Terminal API

### 2.1 终止操作

- (s Stream[T]) CollectToSlice() []T
- (s Stream[T]) ToArray() []T
- (s Stream[T]) ToList() []T
- CollectToMap[T any, K comparable, V any](s Stream[T], keyFn func(T) K, valueFn func(T) V) map[K]V
- (s Stream[T]) Count() int
- (s Stream[T]) Min(less func(a, b T) bool) Optional[T]
- (s Stream[T]) Max(less func(a, b T) bool) Optional[T]
- MinOrdered[T cmp.Ordered](s Stream[T]) Optional[T]
- MaxOrdered[T cmp.Ordered](s Stream[T]) Optional[T]
- (s Stream[T]) AnyMatch(predicate func(T) bool) bool
- (s Stream[T]) AllMatch(predicate func(T) bool) bool
- (s Stream[T]) NoneMatch(predicate func(T) bool) bool
- (s Stream[T]) FindFirst() (T, bool)
- (s Stream[T]) FindAny() (T, bool)
- (s Stream[T]) ForEach(consumer func(T) bool)
- (s Stream[T]) ForEachOrdered(consumer func(T) bool)

### 2.2 迭代与状态

- (s Stream[T]) Seq() iter.Seq[T]
- (s Stream[T]) Iterator() <-chan T
- (s Stream[T]) Unordered() Stream[T]
- (s Stream[T]) Parallel() Stream[T]
- (s Stream[T]) Sequential() Stream[T]
- (s Stream[T]) IsParallel() bool
- (s Stream[T]) OnClose(handler func()) Stream[T]
- (s Stream[T]) Close()

## 3. Collectors API

### 3.1 核心类型与入口

- Collector[T any, A any, R any]
- CollectWith[T any, A any, R any](s Stream[T], c Collector[T, A, R]) R

### 3.2 基础 Collectors

- ToSliceCollector[T any]() Collector[T, []T, []T]
- ToSetCollector[T comparable]() Collector[T, map[T]struct{}, map[T]struct{}]
- CountingCollector[T any]() Collector[T, int64, int64]
- JoiningCollector(delimiter, prefix, suffix string) Collector[string, []string, string]
- MappingCollector[T any, U any, A any, R any](mapper func(T) U, downstream Collector[U, A, R]) Collector[T, A, R]
- FilteringCollector[T any, A any, R any](predicate func(T) bool, downstream Collector[T, A, R]) Collector[T, A, R]
- FlatMappingCollector[T any, U any, A any, R any](mapper func(T) []U, downstream Collector[U, A, R]) Collector[T, A, R]

### 3.3 分组与映射

- GroupingByCollector[T any, K comparable](classifier func(T) K) Collector[T, map[K][]T, map[K][]T]
- GroupingByMappingCollector[T any, K comparable, V any](classifier func(T) K, mapper func(T) V) Collector[T, map[K][]V, map[K][]V]
- GroupingByDownstreamCollector[T any, K comparable, A any, R any](classifier func(T) K, downstream Collector[T, A, R]) Collector[T, map[K]A, map[K]R]
- PartitioningByCollector[T any](predicate func(T) bool) Collector[T, map[bool][]T, map[bool][]T]
- ToMapCollector[T any, K comparable, V any](keyMapper func(T) K, valueMapper func(T) V, mergeFn func(existing, incoming V) V) Collector[T, map[K]V, map[K]V]

### 3.4 归约与组合

- MinByCollector[T any](less func(a, b T) bool) Collector[T, Optional[T], Optional[T]]
- MaxByCollector[T any](less func(a, b T) bool) Collector[T, Optional[T], Optional[T]]
- TeeingCollector[T any, A1 any, R1 any, A2 any, R2 any, R any](left Collector[T, A1, R1], right Collector[T, A2, R2], merger func(R1, R2) R) Collector[T, teeAccumulator[A1, A2], R]
- CollectingAndThen[T any, A any, R any, RR any](downstream Collector[T, A, R], finisher func(R) RR) Collector[T, A, RR]
- ReducingCollector[T any, U any](identity U, mapper func(T) U, op func(U, U) U) Collector[T, U, U]

### 3.5 数值统计

- SummingIntCollector[T any](mapper func(T) int) Collector[T, int64, int64]
- SummingInt64Collector[T any](mapper func(T) int64) Collector[T, int64, int64]
- SummingFloat64Collector[T any](mapper func(T) float64) Collector[T, float64, float64]
- AveragingIntCollector[T any](mapper func(T) int) Collector[T, averagingState, float64]
- AveragingInt64Collector[T any](mapper func(T) int64) Collector[T, averagingState, float64]
- AveragingFloat64Collector[T any](mapper func(T) float64) Collector[T, averagingState, float64]
- IntSummaryStatistics
  - Average() float64
- LongSummaryStatistics
  - Average() float64
- DoubleSummaryStatistics
  - Average() float64
- SummarizingIntCollector[T any](mapper func(T) int) Collector[T, IntSummaryStatistics, IntSummaryStatistics]
- SummarizingInt64Collector[T any](mapper func(T) int64) Collector[T, LongSummaryStatistics, LongSummaryStatistics]
- SummarizingFloat64Collector[T any](mapper func(T) float64) Collector[T, DoubleSummaryStatistics, DoubleSummaryStatistics]

## 4. Spark-like API（本地语义）

### 4.1 变换算子

- Union[T any](left, right Stream[T]) Stream[T]
- Intersection[T comparable](left, right Stream[T]) Stream[T]
- Subtract[T comparable](left, right Stream[T]) Stream[T]
- Cartesian[T any, U any](left Stream[T], right Stream[U]) Stream[Tuple2[T, U]]
- Zip[T any, U any](left Stream[T], right Stream[U]) Stream[Tuple2[T, U]]
- ZipWithIndex[T any](s Stream[T]) Stream[Tuple2[T, int]]
- Sample[T any](s Stream[T], withReplacement bool, fraction float64, seed int64) Stream[T]
- SortBy[T any, K cmp.Ordered](s Stream[T], keyFn func(T) K, ascending bool) Stream[T]
- Glom[T any](s Stream[T], chunkSize int) Stream[[]T]
- (s Stream[T]) Coalesce(_ int) Stream[T]
- (s Stream[T]) Repartition(_ int) Stream[T]

### 4.2 Pair/KV 算子

- MapValues[K comparable, V any, U any](s Stream[Pair[K, V]], mapper func(V) U) Stream[Pair[K, U]]
- FlatMapValues[K comparable, V any, U any](s Stream[Pair[K, V]], mapper func(V) []U) Stream[Pair[K, U]]
- Keys[K comparable, V any](s Stream[Pair[K, V]]) Stream[K]
- Values[K comparable, V any](s Stream[Pair[K, V]]) Stream[V]
- GroupBy[T any, K comparable](s Stream[T], classifier func(T) K) Stream[Pair[K, []T]]
- GroupByKey[K comparable, V any](s Stream[Pair[K, V]]) Stream[Pair[K, []V]]
- ReduceByKey[K comparable, V any](s Stream[Pair[K, V]], reducer func(V, V) V) Stream[Pair[K, V]]
- FoldByKey[K comparable, V any](s Stream[Pair[K, V]], zero V, op func(V, V) V) Stream[Pair[K, V]]
- SortByKey[K cmp.Ordered, V any](s Stream[Pair[K, V]], ascending bool) Stream[Pair[K, V]]
- KeyBy[T any, K comparable](s Stream[T], keyFn func(T) K) Stream[Pair[K, T]]
- Lookup[K comparable, V any](s Stream[Pair[K, V]], key K) []V
- CollectAsMap[K comparable, V any](s Stream[Pair[K, V]]) map[K]V

### 4.3 Join 与聚合

- Join[K comparable, V any, W any](left Stream[Pair[K, V]], right Stream[Pair[K, W]]) Stream[Pair[K, Tuple2[V, W]]]
- LeftOuterJoin[K comparable, V any, W any](left Stream[Pair[K, V]], right Stream[Pair[K, W]]) Stream[Pair[K, Tuple2[V, Optional[W]]]]
- RightOuterJoin[K comparable, V any, W any](left Stream[Pair[K, V]], right Stream[Pair[K, W]]) Stream[Pair[K, Tuple2[Optional[V], W]]]
- FullOuterJoin[K comparable, V any, W any](left Stream[Pair[K, V]], right Stream[Pair[K, W]]) Stream[Pair[K, Tuple2[Optional[V], Optional[W]]]]
- Cogroup[K comparable, V any, W any](left Stream[Pair[K, V]], right Stream[Pair[K, W]]) Stream[Pair[K, Tuple2[[]V, []W]]]
- AggregateByKey[K comparable, V any, U any](s Stream[Pair[K, V]], zero U, seqOp func(U, V) U, combOp func(U, U) U) Stream[Pair[K, U]]
- CombineByKey[K comparable, V any, C any](s Stream[Pair[K, V]], createCombiner func(V) C, mergeValue func(C, V) C, mergeCombiners func(C, C) C) Stream[Pair[K, C]]

### 4.4 Action 算子

- CountByValue[T comparable](s Stream[T]) map[T]int64
- CountByKey[K comparable, V any](s Stream[Pair[K, V]]) map[K]int64
- (s Stream[T]) Take(n int) []T
- (s Stream[T]) First() (T, bool)
- TakeOrdered[T cmp.Ordered](s Stream[T], n int) []T
- Top[T cmp.Ordered](s Stream[T], n int) []T
- MapPartitions[T any, R any](s Stream[T], fn func([]T) []R) Stream[R]
- MapPartitionsWithIndex[T any, R any](s Stream[T], fn func(int, []T) []R) Stream[R]
- ForEachPartition[T any](s Stream[T], consumer func([]T))

## 5. 使用示例

```go
nums := Of(1, 2, 3, 4, 5)
out := Map(nums.Filter(func(v int) bool { return v%2 == 1 }), func(v int) int {
	return v * 10
}).CollectToSlice()
// out == []int{10, 30, 50}
```

```go
pairs := Of(
	Pair[string, int]{Key: "a", Value: 1},
	Pair[string, int]{Key: "a", Value: 2},
	Pair[string, int]{Key: "b", Value: 3},
)
agg := ReduceByKey(pairs, func(x, y int) int { return x + y }).CollectToSlice()
```

## 6. 语义说明

- 所有 Stream 链路默认惰性执行，直到终止操作触发。
- 部分算子会物化数据（如 Sorted、Zip、Cartesian、TakeOrdered、Top）。
- Spark-like API 为单机本地语义，不涉及分布式调度和容错。
- Parallel/Sequential/IsParallel 当前主要用于状态表达，不会自动并行执行计算。

## 7. Go 特有增强 API

### 7.1 Context

- FromContext[T any](ctx context.Context, producer func(context.Context, func(T) bool)) Stream[T]
- (s Stream[T]) WithContext(ctx context.Context) Stream[T]
- (s Stream[T]) TakeUntilContextDone(ctx context.Context) Stream[T]

### 7.2 Channel

- FromChan[T any](ch <-chan T) Stream[T]
- (s Stream[T]) ToChan(buffer int) <-chan T
- MergeChan[T any](chs ...<-chan T) <-chan T
- FanOut[T any](in <-chan T, n int, buffer int) []<-chan T

### 7.3 Error-first

- MapE[T any, R any](s Stream[T], mapper func(T) (R, error)) (Stream[R], error)
- FlatMapE[T any, R any](s Stream[T], mapper func(T) ([]R, error)) (Stream[R], error)
- (s Stream[T]) ForEachE(consumer func(T) error) error
- CollectE[T any, A any](s Stream[T], supplier func() A, accumulator func(*A, T) error, combiner func(*A, A) error) (A, error)
- TapError(err error, handler func(error)) error

### 7.4 并发扩展

- (s Stream[T]) WithWorkerPool(parallelism int) Stream[T]
- ParallelMap[T any, R any](s Stream[T], parallelism int, mapper func(T) R) Stream[R]
- ParallelMapOrdered[T any, R any](s Stream[T], parallelism int, mapper func(T) R) Stream[R]
- ParallelMapUnordered[T any, R any](s Stream[T], parallelism int, mapper func(T) R) Stream[R]
- ParallelFlatMap[T any, R any](s Stream[T], parallelism int, mapper func(T) []R) Stream[R]
- (s Stream[T]) ParallelForEach(parallelism int, consumer func(T))

### 7.5 时间算子

- Window[T any](s Stream[T], size int) Stream[[]T]
- (s Stream[T]) Debounce(interval time.Duration) Stream[T]
- (s Stream[T]) Throttle(interval time.Duration) Stream[T]
- (s Stream[T]) SampleEvery(interval time.Duration) Stream[T]

### 7.6 背压与缓冲

- (s Stream[T]) Buffer(buffer int) Stream[T]
- (s Stream[T]) BlockWhenFull(buffer int) Stream[T]
- (s Stream[T]) DropWhenFull(buffer int) Stream[T]
- (s Stream[T]) LatestOnly() Stream[T]

### 7.7 IO

- FromReaderLines(r io.Reader) Stream[string]
- FromScanner(scanner *bufio.Scanner) Stream[string]
- ToWriter[T any](s Stream[T], w io.Writer, formatter func(T) string) error
- MapBytes(s Stream[[]byte], mapper func([]byte) []byte) Stream[[]byte]

### 7.8 可观测性

- (s Stream[T]) Tap(action func(T)) Stream[T]
- (s Stream[T]) WithMetrics(name string) Stream[T]
- (s Stream[T]) WithTrace(span string) Stream[T]
- SetTraceSink(sink func(TraceEvent))
- GetMetricCount(name string) int64
- ResetMetrics()
