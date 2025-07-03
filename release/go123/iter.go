package main

import "iter"

// https://go.dev/blog/range-functions

// func(yield func() bool)
// func(yield func(v) bool)    => Seq[V any]
// func(yield func(k, v) bool) => Seq[K, V any]

func Filter[V any](f func(V) bool, s iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range s {
			if f(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}
