package channel

import (
	"sync"
	"unsafe"
)

// 三种定义
// 只能接收
// 只能发送
// 既可以接收又可以发送
// ChannelType = ("chan" | "chan<-" | "<-chan") ElementType

/*
	channel 的数据结构

	runtime.hchan
		qcount uint			// 环形队列元素的数量
		dataqsiz uint		// 环形队列的大小
		buf unsafe.Pointer	// 环形队列的指针
		elemsize uint16		// chan中元素的大小
		closed uint32		// 是否已关闭
		elemtype *_type		// 元素的类型
		sendx uint16		// send在buf中的索引
		recvx uint16		// recv在buf中的索引
		recvq waitq			// receiver的等待队列
		sendq waitq			// sender的等待队列
		lock mutex			// 锁, 保护所有的字段
*/

// [qcount] 代表chan中已经接收但是还没有被取走的元素的数量, len(chan) = qcount
// [dataqsiz] 队列的大小. chan使用一个环形队列来存放元素. 环形队列很适合 生产者-消费者的场景.
// [buf] 存放元素的buffer. 在channel创建时,buf就创建好了, 固定的槽位.环形缓冲区,可以重复使用.
// [sendx] 处理返送数据的指针在buf中的位置. 一旦接收到新的数据, 该指针移动到下一个位置.
// [recvx] 处理接收数据的指针在buf中的位置. 一旦发送了数据, 该指针移动到下一个位置.
// [recvq] chan是多生产者,多消费者, 如果消费者因为没有数据可读而被阻塞,那么它就会被加入recvq队列中.
// [sendq] 同上.

const (
	maxAlign  = 8
	hchanSize = unsafe.Sizeof(hchan{}) + uintptr(-int(unsafe.Sizeof(hchan{}))&(maxAlign-1))
)

type hchan struct {
	qcount   uint
	dataqsiz uint
	buf      unsafe.Pointer
	elemsize uint16
	closed   uint32
	elemtype *Type
	sendx    uint16
	recvx    uint16
	recvq    waitq
	sendq    waitq
	lock     sync.Mutex
}

func (c *hchan) raceaddr() unsafe.Pointer { return unsafe.Pointer(&c.buf) }

type waitq struct {
	first *sudog
	last  *sudog
}

// 从队列的前端移除并返回元素
func (q *waitq) dequeue() *sudog {
	for {
		sgp := q.first
		if sgp == nil {
			return nil
		}
		y := sgp.next
		if y == nil {
			q.first = nil
			q.last = nil
		} else {
			y.prev = nil
			q.first = y
			sgp.next = nil
		}
		return sgp
	}
}

type sudog struct {
	next *sudog
	prev *sudog
	g    *g
}

type chantype struct {
	Type
	Elem *Type
}

type Type struct {
	Size_    uintptr
	PtrBytes uintptr // number of (prefix) bytes in the type that can contain pointers
}

func makechan(t *chantype, size int) *hchan {
	elem := t.Elem

	var mem uintptr
	// mem, overflow := math.MulUintptr(elem.Size_, uintptr(size))
	// if overflow || mem > maxAlloc-hchanSize || size < 0 {
	// 	panic(plainError("makechan: size out of range"))
	// }
	var c *hchan
	switch {
	case mem == 0:
		// chan 的容量大小或者元素大小是0, 不必创建buf
		c = (*hchan)(mallocgc(hchanSize, nil, true))
		c.buf = c.raceaddr()
	case elem.PtrBytes == 0:
		// 元素不是指针, 分配一块连续的内存给hchan数据结构和buf
		c = (*hchan)(mallocgc(hchanSize+mem, nil, true))
		c.buf = add(unsafe.Pointer(c), hchanSize)
	default:
		// 元素包含指针,单独分配buf
		c = new(hchan)
		c.buf = mallocgc(mem, elem, true)
	}
	// 元素的大小、元素的类型、chan的容量都被记录下来
	c.elemsize = uint16(elem.Size_)
	c.elemtype = elem
	c.dataqsiz = uint(size)

	return c
}

func throw(s string)                                                 {}
func mallocgc(size uintptr, typ *Type, needzero bool) unsafe.Pointer { return nil }

func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

func getcallerpc() uintptr {
	return 0
}
func lock(l *sync.Mutex) {
	// lockWithRank(l, getLockRank(l))
}

func unlock(l *sync.Mutex) {
	// unlockWithRank(l, getLockRank(l))
}