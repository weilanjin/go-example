package main

import (
	"runtime"
	"sync"
)

func resolveParallelism(parallelism int, taskCount int) int {
	if taskCount <= 0 {
		return 1
	}
	if parallelism <= 0 {
		parallelism = runtime.GOMAXPROCS(0)
	}
	if parallelism < 1 {
		parallelism = 1
	}
	if parallelism > taskCount {
		parallelism = taskCount
	}
	return parallelism
}

func cloneStreamState[T any, R any](from Stream[T], seq func(func(R) bool), sizeHint int) Stream[R] {
	res := makeStream(seq, sizeHint)
	res.parallel = from.parallel
	res.workerPool = from.workerPool
	res.unordered = from.unordered
	if len(from.onClose) > 0 {
		res.onClose = append([]func(){}, from.onClose...)
	}
	return res
}

func parallelismForStream[T any](s Stream[T], parallelism int, taskCount int) int {
	if parallelism > 0 {
		return resolveParallelism(parallelism, taskCount)
	}
	if s.workerPool > 0 {
		return resolveParallelism(s.workerPool, taskCount)
	}
	return resolveParallelism(parallelism, taskCount)
}

// MapPar maps elements in parallel while preserving input order.
// MapPar 并行映射元素并保持输入顺序。
func MapPar[T any, R any](s Stream[T], parallelism int, mapper func(T) R) Stream[R] {
	if mapper == nil {
		return Empty[R]()
	}

	seq := func(yield func(R) bool) {
		items := s.Slice()
		n := len(items)
		if n == 0 {
			return
		}

		workerCount := parallelismForStream(s, parallelism, n)
		jobs := make(chan int, workerCount*2)
		out := make([]R, n)

		var wg sync.WaitGroup
		wg.Add(workerCount)
		for w := 0; w < workerCount; w++ {
			go func() {
				defer wg.Done()
				for idx := range jobs {
					out[idx] = mapper(items[idx])
				}
			}()
		}

		for i := 0; i < n; i++ {
			jobs <- i
		}
		close(jobs)
		wg.Wait()

		for i := 0; i < n; i++ {
			if !yield(out[i]) {
				return
			}
		}
	}

	return cloneStreamState(s, seq, s.sizeHint)
}

// MapParUnordered maps elements in parallel and emits as tasks complete.
// MapParUnordered 并行映射元素并按任务完成顺序输出。
func MapParUnordered[T any, R any](s Stream[T], parallelism int, mapper func(T) R) Stream[R] {
	if mapper == nil {
		return Empty[R]()
	}

	seq := func(yield func(R) bool) {
		items := s.Slice()
		n := len(items)
		if n == 0 {
			return
		}

		workerCount := parallelismForStream(s, parallelism, n)
		jobs := make(chan int, workerCount*2)
		results := make(chan R, workerCount*2)

		var wg sync.WaitGroup
		wg.Add(workerCount)
		for w := 0; w < workerCount; w++ {
			go func() {
				defer wg.Done()
				for idx := range jobs {
					results <- mapper(items[idx])
				}
			}()
		}

		for i := 0; i < n; i++ {
			jobs <- i
		}
		close(jobs)

		go func() {
			wg.Wait()
			close(results)
		}()

		stopped := false
		for r := range results {
			if stopped {
				continue
			}
			if !yield(r) {
				stopped = true
			}
		}
	}

	return cloneStreamState(s, seq, s.sizeHint)
}

// FlatMapPar maps each element to a slice in parallel and flattens results.
// FlatMapPar 并行将元素映射为切片并扁平化结果。
func FlatMapPar[T any, R any](s Stream[T], parallelism int, mapper func(T) []R) Stream[R] {
	if mapper == nil {
		return Empty[R]()
	}

	mapped := MapPar(s, parallelism, mapper).Slice()
	total := 0
	for i := range mapped {
		total += len(mapped[i])
	}
	out := make([]R, 0, total)
	for i := range mapped {
		out = append(out, mapped[i]...)
	}
	return FromSlice(out)
}

// EachPar applies consumer to each element in parallel.
// EachPar 对每个元素并行执行 consumer。
func (s Stream[T]) EachPar(parallelism int, consumer func(T)) {
	if consumer == nil {
		return
	}

	items := s.Slice()
	n := len(items)
	if n == 0 {
		return
	}

	workerCount := parallelismForStream(s, parallelism, n)
	jobs := make(chan int, workerCount*2)

	var wg sync.WaitGroup
	wg.Add(workerCount)
	for w := 0; w < workerCount; w++ {
		go func() {
			defer wg.Done()
			for idx := range jobs {
				consumer(items[idx])
			}
		}()
	}

	for i := 0; i < n; i++ {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
}
