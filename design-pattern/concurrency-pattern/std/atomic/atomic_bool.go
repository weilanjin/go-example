package atomic

import "sync/atomic"

type Bool struct {
	_ noCopy
	v uint32 // 底层使用了 uint32
}

func (b *Bool) Load() bool {
	return atomic.LoadUint32(&b.v) != 0
}

func (b *Bool) Store(val bool) {
	atomic.StoreUint32(&b.v, b32(val))
}

func (b *Bool) Swap(new bool) bool {
	return atomic.SwapUint32(&b.v, b32(new)) != 0
}

func (b *Bool) CompareAndSwap(old, new bool) bool {
	return atomic.CompareAndSwapUint32(&b.v, b32(old), b32(new))
}

func b32(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}
