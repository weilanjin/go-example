package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type TraceEvent struct {
	Span     string
	Count    int
	Duration time.Duration
}

var traceSink atomic.Value
var metricCounters sync.Map

func init() {
	traceSink.Store(func(TraceEvent) {})
}

// SetTraceSink sets a global trace sink callback.
// SetTraceSink 设置全局 trace 回调。
func SetTraceSink(sink func(TraceEvent)) {
	if sink == nil {
		traceSink.Store(func(TraceEvent) {})
		return
	}
	traceSink.Store(sink)
}

// GetMetricCount returns count by metric name.
// GetMetricCount 按指标名返回计数。
func GetMetricCount(name string) int64 {
	if name == "" {
		return 0
	}
	v, ok := metricCounters.Load(name)
	if !ok {
		return 0
	}
	ptr, ok := v.(*int64)
	if !ok {
		return 0
	}
	return atomic.LoadInt64(ptr)
}

// ResetMetrics clears all in-memory metric counters.
// ResetMetrics 清空内存中的全部指标计数。
func ResetMetrics() {
	metricCounters.Range(func(key any, _ any) bool {
		metricCounters.Delete(key)
		return true
	})
}

func metricPtr(name string) *int64 {
	if name == "" {
		return nil
	}
	if ptr, ok := metricCounters.Load(name); ok {
		return ptr.(*int64)
	}
	zero := new(int64)
	actual, _ := metricCounters.LoadOrStore(name, zero)
	return actual.(*int64)
}

// TapError invokes handler when err is non-nil and returns err unchanged.
// TapError 在 err 非空时执行处理函数并原样返回 err。
func TapError(err error, handler func(error)) error {
	if err != nil && handler != nil {
		handler(err)
	}
	return err
}

// Tap applies side effects without changing stream elements.
// Tap 在不改变元素的前提下执行副作用。
func (s Stream[T]) Tap(action func(T)) Stream[T] {
	return s.Peek(action)
}

// WithMetrics counts passed-through elements under given metric name.
// WithMetrics 对经过的元素按指标名计数。
func (s Stream[T]) WithMetrics(name string) Stream[T] {
	ptr := metricPtr(name)
	if ptr == nil {
		return s
	}
	seq := func(yield func(T) bool) {
		s.seq(func(v T) bool {
			atomic.AddInt64(ptr, 1)
			return yield(v)
		})
	}
	return s.derive(seq, s.sizeHint)
}

// WithTrace emits one trace event after stream traversal finishes.
// WithTrace 在遍历结束后输出一条 trace 事件。
func (s Stream[T]) WithTrace(span string) Stream[T] {
	if span == "" {
		return s
	}
	seq := func(yield func(T) bool) {
		start := time.Now()
		count := 0
		s.seq(func(v T) bool {
			count++
			return yield(v)
		})
		sink := traceSink.Load().(func(TraceEvent))
		sink(TraceEvent{Span: span, Count: count, Duration: time.Since(start)})
	}
	return s.derive(seq, s.sizeHint)
}

// FromContext builds a stream driven by a context-aware producer.
// FromContext 基于感知 context 的生产函数创建流。
func FromContext[T any](ctx context.Context, producer func(context.Context, func(T) bool)) Stream[T] {
	if ctx == nil || producer == nil {
		return Empty[T]()
	}
	seq := func(yield func(T) bool) {
		stopped := false
		producer(ctx, func(v T) bool {
			if stopped {
				return false
			}
			select {
			case <-ctx.Done():
				stopped = true
				return false
			default:
			}
			if !yield(v) {
				stopped = true
				return false
			}
			return true
		})
	}
	return makeStream(seq, -1)
}

// FromCtx is a short alias of FromContext.
// FromCtx 是 FromContext 的简写别名。
func FromCtx[T any](ctx context.Context, producer func(context.Context, func(T) bool)) Stream[T] {
	return FromContext(ctx, producer)
}

// WithContext stops iteration when ctx is done.
// WithContext 在 ctx 结束时停止迭代。
func (s Stream[T]) WithContext(ctx context.Context) Stream[T] {
	if ctx == nil {
		return s
	}
	seq := func(yield func(T) bool) {
		s.seq(func(v T) bool {
			select {
			case <-ctx.Done():
				return false
			default:
			}
			return yield(v)
		})
	}
	return s.derive(seq, -1)
}

// WithCtx is a short alias of WithContext.
// WithCtx 是 WithContext 的简写别名。
func (s Stream[T]) WithCtx(ctx context.Context) Stream[T] {
	return s.WithContext(ctx)
}

// TakeUntilContextDone keeps elements until ctx is done.
// TakeUntilContextDone 持续输出直到 ctx 结束。
func (s Stream[T]) TakeUntilContextDone(ctx context.Context) Stream[T] {
	return s.WithContext(ctx)
}

// FromChan builds a stream from a read-only channel.
// FromChan 从只读 channel 创建流。
func FromChan[T any](ch <-chan T) Stream[T] {
	if ch == nil {
		return Empty[T]()
	}
	seq := func(yield func(T) bool) {
		for v := range ch {
			if !yield(v) {
				return
			}
		}
	}
	return makeStream(seq, -1)
}

// FromChannel is an alias of FromChan.
// FromChannel 是 FromChan 的别名。
func FromChannel[T any](ch <-chan T) Stream[T] {
	return FromChan(ch)
}

// ToChan converts stream to channel with default background context.
// ToChan 使用默认后台 context 将流转换为 channel。
func (s Stream[T]) ToChan(buffer int) <-chan T {
	return s.ToChanCtx(context.Background(), buffer)
}

// ToChanCtx converts stream to channel and supports cancellation by ctx.
// ToChanCtx 将流转换为 channel，并支持 ctx 取消。
func (s Stream[T]) ToChanCtx(ctx context.Context, buffer int) <-chan T {
	if ctx == nil {
		ctx = context.Background()
	}
	if buffer < 0 {
		buffer = 0
	}
	ch := make(chan T, buffer)
	go func() {
		defer close(ch)
		s.seq(func(v T) bool {
			select {
			case <-ctx.Done():
				return false
			case ch <- v:
				return true
			}
		})
	}()
	return ch
}

// MergeChan merges multiple channels into one output channel.
// MergeChan 将多个 channel 合并为一个输出 channel。
func MergeChan[T any](chs ...<-chan T) <-chan T {
	return MergeChanCtx(context.Background(), chs...)
}

// MergeChanCtx merges channels with context cancellation support.
// MergeChanCtx 在支持 context 取消的情况下合并多个 channel。
func MergeChanCtx[T any](ctx context.Context, chs ...<-chan T) <-chan T {
	if ctx == nil {
		ctx = context.Background()
	}
	out := make(chan T)
	if len(chs) == 0 {
		close(out)
		return out
	}
	var wg sync.WaitGroup
	wg.Add(len(chs))
	for i := range chs {
		ch := chs[i]
		go func() {
			defer wg.Done()
			for v := range ch {
				select {
				case <-ctx.Done():
					return
				case out <- v:
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// MergeChannels is an alias of MergeChan.
// MergeChannels 是 MergeChan 的别名。
func MergeChannels[T any](chs ...<-chan T) <-chan T {
	return MergeChan(chs...)
}

// FanOut broadcasts input channel with blocking semantics by default.
// FanOut 默认以阻塞语义广播输入 channel。
func FanOut[T any](in <-chan T, n int, buffer int) []<-chan T {
	return FanOutBlock(in, n, buffer)
}

// FanOutBlock broadcasts values and blocks on slow consumers.
// FanOutBlock 广播值，并在慢消费者处阻塞。
func FanOutBlock[T any](in <-chan T, n int, buffer int) []<-chan T {
	return FanOutBlockCtx(context.Background(), in, n, buffer)
}

// FanOutDrop broadcasts values and drops when consumer buffer is full.
// FanOutDrop 广播值，并在消费者缓冲满时丢弃。
func FanOutDrop[T any](in <-chan T, n int, buffer int) []<-chan T {
	return FanOutDropCtx(context.Background(), in, n, buffer)
}

// FanOutBlockCtx broadcasts with blocking behavior and ctx cancellation.
// FanOutBlockCtx 使用阻塞广播语义，并支持 ctx 取消。
func FanOutBlockCtx[T any](ctx context.Context, in <-chan T, n int, buffer int) []<-chan T {
	if ctx == nil {
		ctx = context.Background()
	}
	if n <= 0 {
		return []<-chan T{}
	}
	if buffer < 0 {
		buffer = 0
	}
	outs := make([]chan T, n)
	ro := make([]<-chan T, n)
	for i := 0; i < n; i++ {
		outs[i] = make(chan T, buffer)
		ro[i] = outs[i]
	}
	go func() {
		defer func() {
			for i := range outs {
				close(outs[i])
			}
		}()
		for {
			var (
				v  T
				ok bool
			)
			select {
			case <-ctx.Done():
				return
			case v, ok = <-in:
				if !ok {
					return
				}
			}
			for i := range outs {
				select {
				case <-ctx.Done():
					return
				case outs[i] <- v:
				}
			}
		}
	}()
	return ro
}

// FanOutDropCtx broadcasts with drop behavior and ctx cancellation.
// FanOutDropCtx 使用丢弃广播语义，并支持 ctx 取消。
func FanOutDropCtx[T any](ctx context.Context, in <-chan T, n int, buffer int) []<-chan T {
	if ctx == nil {
		ctx = context.Background()
	}
	if n <= 0 {
		return []<-chan T{}
	}
	if buffer < 0 {
		buffer = 0
	}
	outs := make([]chan T, n)
	ro := make([]<-chan T, n)
	for i := 0; i < n; i++ {
		outs[i] = make(chan T, buffer)
		ro[i] = outs[i]
	}
	go func() {
		defer func() {
			for i := range outs {
				close(outs[i])
			}
		}()
		for {
			var (
				v  T
				ok bool
			)
			select {
			case <-ctx.Done():
				return
			case v, ok = <-in:
				if !ok {
					return
				}
			}
			for i := range outs {
				select {
				case <-ctx.Done():
					return
				case outs[i] <- v:
				default:
				}
			}
		}
	}()
	return ro
}

// MapE maps with error and returns first error immediately.
// MapE 带错误映射，并在首个错误出现时立即返回。
func MapE[T any, R any](s Stream[T], mapper func(T) (R, error)) (Stream[R], error) {
	if mapper == nil {
		return Empty[R](), nil
	}
	items := s.Slice()
	out := make([]R, 0, len(items))
	for i := range items {
		v, err := mapper(items[i])
		if err != nil {
			return Empty[R](), err
		}
		out = append(out, v)
	}
	return FromSlice(out), nil
}

// TryMap is an alias of MapE.
// TryMap 是 MapE 的别名。
func TryMap[T any, R any](s Stream[T], mapper func(T) (R, error)) (Stream[R], error) {
	return MapE(s, mapper)
}

// FlatMapE flat-maps with error and returns first error immediately.
// FlatMapE 带错误扁平映射，并在首个错误出现时立即返回。
func FlatMapE[T any, R any](s Stream[T], mapper func(T) ([]R, error)) (Stream[R], error) {
	if mapper == nil {
		return Empty[R](), nil
	}
	items := s.Slice()
	out := make([]R, 0, len(items))
	for i := range items {
		arr, err := mapper(items[i])
		if err != nil {
			return Empty[R](), err
		}
		out = append(out, arr...)
	}
	return FromSlice(out), nil
}

// TryFlatMap is an alias of FlatMapE.
// TryFlatMap 是 FlatMapE 的别名。
func TryFlatMap[T any, R any](s Stream[T], mapper func(T) ([]R, error)) (Stream[R], error) {
	return FlatMapE(s, mapper)
}

// ForEachE iterates elements until consumer returns an error.
// ForEachE 遍历元素，直到 consumer 返回错误。
func (s Stream[T]) ForEachE(consumer func(T) error) error {
	if consumer == nil {
		return nil
	}
	var firstErr error
	s.seq(func(v T) bool {
		err := consumer(v)
		if err != nil {
			firstErr = err
			return false
		}
		return true
	})
	return firstErr
}

// CollectE collects with error-aware accumulator/combiner.
// CollectE 使用支持错误的 accumulator/combiner 收集数据。
func CollectE[T any, A any](s Stream[T], supplier func() A, accumulator func(*A, T) error, combiner func(*A, A) error) (A, error) {
	var zero A
	if supplier == nil {
		return zero, nil
	}
	acc := supplier()
	if accumulator != nil {
		for _, v := range s.Slice() {
			if err := accumulator(&acc, v); err != nil {
				return zero, err
			}
		}
	}
	if combiner != nil {
		empty := supplier()
		if err := combiner(&acc, empty); err != nil {
			return zero, err
		}
	}
	return acc, nil
}

// TryCollect is an alias of CollectE.
// TryCollect 是 CollectE 的别名。
func TryCollect[T any, A any](s Stream[T], supplier func() A, accumulator func(*A, T) error, combiner func(*A, A) error) (A, error) {
	return CollectE(s, supplier, accumulator, combiner)
}

// Window groups stream into fixed-size chunks.
// Window 按固定大小将流分组。
func Window[T any](s Stream[T], size int) Stream[[]T] {
	return Glom(s, size)
}

// Batch is an alias of Window.
// Batch 是 Window 的别名。
func Batch[T any](s Stream[T], size int) Stream[[]T] {
	return Window(s, size)
}

// Buffer creates a buffered bridge stream with background context.
// Buffer 使用后台 context 创建带缓冲的桥接流。
func (s Stream[T]) Buffer(buffer int) Stream[T] {
	return s.BufferCtx(context.Background(), buffer)
}

// BufferCtx creates a buffered bridge stream with context cancellation.
// BufferCtx 使用可取消 context 创建带缓冲的桥接流。
func (s Stream[T]) BufferCtx(ctx context.Context, buffer int) Stream[T] {
	if ctx == nil {
		ctx = context.Background()
	}
	if buffer < 0 {
		buffer = 0
	}
	ch := make(chan T, buffer)
	go func() {
		defer close(ch)
		s.seq(func(v T) bool {
			select {
			case <-ctx.Done():
				return false
			case ch <- v:
				return true
			}
		})
	}()
	return FromChan(ch)
}

// BlockWhenFull blocks producer when buffer is full.
// BlockWhenFull 在缓冲满时阻塞生产者。
func (s Stream[T]) BlockWhenFull(buffer int) Stream[T] {
	return s.Buffer(buffer)
}

// BlockWhenFullCtx blocks producer with context cancellation support.
// BlockWhenFullCtx 在支持 context 取消的情况下阻塞生产者。
func (s Stream[T]) BlockWhenFullCtx(ctx context.Context, buffer int) Stream[T] {
	return s.BufferCtx(ctx, buffer)
}

// DropWhenFull drops elements when buffer is full.
// DropWhenFull 在缓冲满时丢弃元素。
func (s Stream[T]) DropWhenFull(buffer int) Stream[T] {
	return s.DropWhenFullCtx(context.Background(), buffer)
}

// DropWhenFullCtx drops elements with context cancellation support.
// DropWhenFullCtx 在支持 context 取消的情况下丢弃元素。
func (s Stream[T]) DropWhenFullCtx(ctx context.Context, buffer int) Stream[T] {
	if ctx == nil {
		ctx = context.Background()
	}
	if buffer <= 0 {
		buffer = 1
	}
	ch := make(chan T, buffer)
	go func() {
		defer close(ch)
		s.seq(func(v T) bool {
			select {
			case <-ctx.Done():
				return false
			case ch <- v:
			default:
			}
			return true
		})
	}()
	return FromChan(ch)
}

// LatestOnly keeps only the latest element under pressure.
// LatestOnly 在压力场景下仅保留最新元素。
func (s Stream[T]) LatestOnly() Stream[T] {
	return s.LatestOnlyCtx(context.Background())
}

// LatestOnlyCtx keeps only latest element with context cancellation support.
// LatestOnlyCtx 在支持 context 取消的情况下仅保留最新元素。
func (s Stream[T]) LatestOnlyCtx(ctx context.Context) Stream[T] {
	if ctx == nil {
		ctx = context.Background()
	}
	ch := make(chan T, 1)
	go func() {
		defer close(ch)
		s.seq(func(v T) bool {
			select {
			case <-ctx.Done():
				return false
			case ch <- v:
			default:
				<-ch
				select {
				case <-ctx.Done():
					return false
				case ch <- v:
				}
			}
			return true
		})
	}()
	return FromChan(ch)
}

// Throttle limits output rate by interval.
// Throttle 按时间间隔限制输出速率。
func (s Stream[T]) Throttle(interval time.Duration) Stream[T] {
	if interval <= 0 {
		return s
	}
	seq := func(yield func(T) bool) {
		next := time.Time{}
		s.seq(func(v T) bool {
			now := time.Now()
			if now.Before(next) {
				return true
			}
			next = now.Add(interval)
			return yield(v)
		})
	}
	return s.derive(seq, -1)
}

// Debounce emits the latest value after quiet interval.
// Debounce 在静默间隔后输出最新值。
func (s Stream[T]) Debounce(interval time.Duration) Stream[T] {
	if interval <= 0 {
		return s
	}
	seq := func(yield func(T) bool) {
		in := s.ToChan(0)
		var (
			last  T
			has   bool
			timer *time.Timer
		)
		flush := func() bool {
			if !has {
				return true
			}
			has = false
			return yield(last)
		}
		for {
			if timer == nil {
				v, ok := <-in
				if !ok {
					return
				}
				last = v
				has = true
				timer = time.NewTimer(interval)
				continue
			}
			select {
			case v, ok := <-in:
				if !ok {
					if !timer.Stop() {
						select {
						case <-timer.C:
						default:
						}
					}
					_ = flush()
					return
				}
				last = v
				has = true
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				timer.Reset(interval)
			case <-timer.C:
				if !flush() {
					return
				}
				timer = nil
			}
		}
	}
	return s.derive(seq, -1)
}

// SampleEvery samples the latest value on each interval tick.
// SampleEvery 在每个时间间隔采样最新值。
func (s Stream[T]) SampleEvery(interval time.Duration) Stream[T] {
	if interval <= 0 {
		return s
	}
	seq := func(yield func(T) bool) {
		in := s.ToChan(0)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		var (
			latest T
			has    bool
		)
		for {
			select {
			case v, ok := <-in:
				if !ok {
					if has {
						_ = yield(latest)
					}
					return
				}
				latest = v
				has = true
			case <-ticker.C:
				if has {
					if !yield(latest) {
						return
					}
					has = false
				}
			}
		}
	}
	return s.derive(seq, -1)
}

// FromScanner builds a line stream from scanner text tokens.
// FromScanner 基于 scanner 文本 token 创建行流。
func FromScanner(scanner *bufio.Scanner) Stream[string] {
	if scanner == nil {
		return Empty[string]()
	}
	seq := func(yield func(string) bool) {
		for scanner.Scan() {
			if !yield(scanner.Text()) {
				return
			}
		}
	}
	return makeStream(seq, -1)
}

// FromReaderLines builds a line stream from io.Reader.
// FromReaderLines 基于 io.Reader 创建行流。
func FromReaderLines(r io.Reader) Stream[string] {
	if r == nil {
		return Empty[string]()
	}
	return FromScanner(bufio.NewScanner(r))
}

// ToWriter writes stream elements to io.Writer using formatter.
// ToWriter 使用 formatter 将流元素写入 io.Writer。
func ToWriter[T any](s Stream[T], w io.Writer, formatter func(T) string) error {
	if w == nil {
		return fmt.Errorf("nil writer")
	}
	if formatter == nil {
		formatter = func(v T) string { return fmt.Sprint(v) }
	}
	for _, v := range s.Slice() {
		if _, err := io.WriteString(w, formatter(v)); err != nil {
			return err
		}
	}
	return nil
}

// MapBytes maps byte-slice elements and returns copied mapped bytes.
// MapBytes 映射字节切片元素，并返回复制后的映射结果。
func MapBytes(s Stream[[]byte], mapper func([]byte) []byte) Stream[[]byte] {
	if mapper == nil {
		return Empty[[]byte]()
	}
	return Map(s, func(v []byte) []byte {
		mapped := mapper(v)
		if mapped == nil {
			return nil
		}
		out := make([]byte, len(mapped))
		copy(out, mapped)
		return out
	})
}

// defaultParallelism returns a safe default worker count.
// defaultParallelism 返回安全的默认 worker 数量。
func defaultParallelism() int {
	n := runtime.GOMAXPROCS(0)
	if n < 1 {
		return 1
	}
	return n
}
