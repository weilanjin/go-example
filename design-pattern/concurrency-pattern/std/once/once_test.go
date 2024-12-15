package once

import (
	"math/big"
	"sync"
	"testing"
)

// example 1

var threeOnce struct { // 封装成单例
	sync.Once
	v *big.Float
}

func three() *big.Float {
	threeOnce.Do(func() {
		threeOnce.v = big.NewFloat(3.0)
	})
	return threeOnce.v
}

// example 2

func OnceFn[T any](f func() T) func() T {
	var once sync.Once
	var t T
	fn := func() {
		t = f()
	}
	return func() T {
		once.Do(fn)
		return t
	}
}

// deadlock

func TestDeadlock(t *testing.T) {
	var once sync.Once
	once.Do(func() {
		once.Do(func() {}) // deadlock
	})
}
