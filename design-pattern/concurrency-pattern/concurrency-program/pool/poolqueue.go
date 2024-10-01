package pool

import (
	"sync/atomic"
	"unsafe"
)

type poolDequeue struct {
	headTail atomic.Uint64
	vals     []eface
}

type eface struct {
	typ, val unsafe.Pointer
}

const dequeueBits = 32
const dequeueLimit = (1 << dequeueBits) / 4

type poolChain struct {
	// head is the poolDequeue to push to. This is only accessed
	// by the producer, so doesn't need to be synchronized.
	head *poolChainElt

	// tail is the poolDequeue to popTail from. This is accessed
	// by consumers, so reads and writes must be atomic.
	tail *poolChainElt
}

func (p *poolChain) pushHead(val any) bool {
	return false
}

func (p *poolChain) popHead() (any, bool) {
	return nil, false
}

func (p *poolChain) pushTail(val any) bool {
	return false
}

func (p *poolChain) popTail() (any, bool) {
	return nil, false
}

type poolChainElt struct {
	poolDequeue
	next, prev *poolChainElt
}
