package channel

type Stream[T any] struct {
	ch   chan T
	done chan struct{}
}

func AsStream[T any](done chan struct{}, values ...T) Stream[T] {
	s := Stream[T]{
		ch:   make(chan T),
		done: done,
	}
	go func() {
		defer close(s.ch)
		for _, v := range values {
			select {
			case <-done:
				return
			case s.ch <- v:
			}
		}
	}()
	return s
}

func (s Stream[T]) TakeN(n int) <-chan T {
	res := make(chan T)
	go func() {
		defer close(res)
		for range n {
			res <- <-s.ch
		}
		s.done <- struct{}{}
	}()
	return res
}