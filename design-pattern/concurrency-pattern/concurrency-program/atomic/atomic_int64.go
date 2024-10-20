package atomic

import "sync/atomic"

// 辅助 vet 等 int 工具检查用的, 检查 int64有没有被复制使用,它不占用额外的字节
// Note that it must not be embedded, due to the Lock and Unlock methods.
type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

// align64 不占用额外字节
// 告诉编译器要64bit对齐, 因为对 int64 的原子操作必须要求64对齐
// 32bit的架构中, 如果不对齐,则可能会导致panic
//
// align64 may be added to structs that must be 64-bit aligned.
// This struct is recognized by a special case in the compiler
// and will not work if copied to any other package.
type align64 struct{}

type Uint64 struct {
	_ noCopy
	_ align64 // 对齐标志
	v uint64
}

func (x *Uint64) Load() uint64 {
	return atomic.LoadUint64(&x.v)
}

func (x *Uint64) Store(i uint64) {
	atomic.StoreUint64(&x.v, i)
}

func (x *Uint64) Swap(i uint64) uint64 {
	return atomic.SwapUint64(&x.v, i)
}

func (x *Uint64) CompareAndSwap(old, new uint64) bool {
	return atomic.CompareAndSwapUint64(&x.v, old, new)
}

func (x *Uint64) Add(delta uint64) uint64 {
	return atomic.AddUint64(&x.v, delta)
}
