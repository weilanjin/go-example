package channel

func mapChan[T, K any](in <-chan T, fn func(T) K) <-chan K {
	out := make(chan K)
	if in == nil {
		close(out)
		return out
	}
	go func() {
		defer close(out)
		for v := range in {
			out <- fn(v)
		}
	}()
	return out
}

func reduce[T, K any](in <-chan T, fn func(K, T) K, initial K) K {
	out := initial
	if in == nil {
		return out
	}
	for v := range in {
		out = fn(out, v)
	}
	return out
}