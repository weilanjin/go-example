package atomic

import (
	"log"
	"sync/atomic"
	"testing"
)

func TestAtomic(t *testing.T) {

	// AddUxx 无符号类型操作

	var x uint64 = 0
	newX := atomic.AddUint64(&x, 100) // newX = 100
	log.Print(newX)

	newX = atomic.AddUint64(&x, ^uint64(0)) // newX = 99
	log.Println(newX)

	newX = atomic.AddUint64(&x, ^uint64(10-1)) // x == 89
	log.Println(newX)

	// CompareAndSwapUxx

	x = 0
	ok := atomic.CompareAndSwapUint64(&x, 0, 100) // true, x == 100
	log.Println(ok, x)

	ok = atomic.CompareAndSwapUint64(&x, 0, 100) // false, x 的 旧值不是0
	log.Println(ok, x)

	// SwapXxx

	x = 0
	old := atomic.SwapUint64(&x, 100) // old == 0
	log.Println(old)

	old = atomic.SwapUint64(&x, 100) // old == 100
	log.Println(old)

	// LoadXxx

	x = 0
	v := atomic.LoadUint64(&x) // v == 0
	log.Println(v)

	x = 100
	v = atomic.LoadUint64(&x) // v == 100
	log.Println(v)

	// StoreXxx
	x = 0
	atomic.StoreUint64(&x, 100) // x == 100
	log.Println(x)

}
