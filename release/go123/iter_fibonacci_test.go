package main

import (
	"fmt"
	"iter"
	"testing"
)

func fib0() iter.Seq[int] {
	return func(yield func(int) bool) {
		for a, b := 0, 1; yield(a); a, b = b, a+b {
		}
	}
}

func TestFib0(t *testing.T) {
	for n := range fib0() {
		if n > 100 {
			break
		}
		fmt.Printf("%d ", n)
	}
}
