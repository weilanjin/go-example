package pool

import (
	"math/bits"
	"sync"
)

// buffer 做回收时,要检查buffer 的大小
// 如果buffer太大了,就不会回收,会照成内存泄露
// e.Cap > 64<<10 //64kb
//
// 为了节省内存空间. buffer 采用分级设计
// 常用的第三方库
// bytebufferpool: 这个是 fasthttp 作者 valyala 提供的一个buffer 池
// oxtoacart/bpool:
//  bpool: 是基于 channel实现, 限制 Pool 容量, 一旦返回时数量大于它的阈值,就会自动丢弃
// 1.bpool.BufferPool 提供固定数量的buffer池
// 2.bpool.BytesPool  提供固定数量的byte slice 池
// 3.bpool.SizedBufferPool

type BytePool struct {
	pools            []*sizedPool
	minSize, maxSize int
}

func (p *BytePool) findPool(size int) *sizedPool { // 选择指定大小的池子
	if size > p.maxSize {
		return nil
	}
	div, rem := bits.Div64(0, uint64(size), uint64(p.minSize))
	idx := bits.Len64(div)
	if rem == 0 && div != 0 && (div&(div-1)) == 0 {
		idx--
	}
	return p.pools[idx]
}

func (p *BytePool) Get(size int) *[]byte { // 获取对象
	sp := p.findPool(size) // 先找到对应的池子
	if sp == nil {
		return makeSlicePointer(size)
	}
	buf := sp.pool.Get().(*[]byte)
	*buf = (*buf)[:size]
	return buf
}

func (p *BytePool) Put(buf *[]byte) { // 将对象放回池子
	sp := p.findPool(cap(*buf)) // 找到对应的池子
	if sp == nil {
		return
	}
	*buf = (*buf)[:cap(*buf)]
	sp.pool.Put(buf)
}

type sizedPool struct { // sizePool 包含对应的大小和sync.Pool
	pool sync.Pool
	size int
}

func newSizedPool(size int) *sizedPool {
	return &sizedPool{
		pool: sync.Pool{
			New: func() any {
				return makeSlicePointer(size)
			},
		},
		size: size,
	}
}

func makeSlicePointer(size int) *[]byte {
	b := make([]byte, 0, size)
	return &b
}
