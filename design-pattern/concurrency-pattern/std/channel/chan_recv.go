package channel

import (
	"unsafe"
)

func chanrecv1(c *hchan, elem unsafe.Pointer) {
	chanrecv(c, elem, true)
}

func chanrecv2(c *hchan, elem unsafe.Pointer) (received bool) {
	_, received = chanrecv(c, elem, false)
	return
}

func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) {
	// 1.chan 的值为 nil
	if c == nil {
		if !block {
			return
		}
		// gopark(nil, nil, waitReasonChanReceiveNilChan, traceEvGoStop, 2)
		throw("unreachable")
	}
	// 2.block 参数的值为false且c 为空
	if !block && empty(c) {
		// ....
	}
	lock(&c.lock)
	// 3.chan 已经被关闭, 且为空, 已经没有数据了
	if c.closed != 0 && c.qcount == 0 {
		unlock(&c.lock)
		return true, false
	}
	// 4.如果 sendq队列中有等待发送的sender
	if sg := c.sendq.dequeue(); sg != nil {
		recv(c, sg, ep, func() { unlock(&c.lock) }, 3)
		return true, true
	}
	// 5.没有等待的sender, buf中有数据
	if c.qcount > 0 {
		/*
			qp := chanbuf(c, c.recvq.tail)
			if ep != nil {
				typedmemmove(c, elemtype, ep, qp)
			}
			typedmemclr(c.elemtype, qp)
			c.recvx++
			if c.recvx == c.dataqsiz {
				c.recvx = 0
			}
			c.qcount--
			unlock(&c.lock)
			return true, true
		*/
	}
	if !block {
		unlock(&c.lock)
		return false, false
	}
	// 6. buf 中没有数据直接阻塞
	// ....
	return
}

func recv(c *hchan, sg *sudog, ep unsafe.Pointer, unlockf func(), skip int) {}

func empty(c *hchan) bool {
	//// c.dataqsiz is immutable.
	//if c.dataqsiz == 0 {
	//	return atomic.Loadp(unsafe.Pointer(&c.sendq.first)) == nil
	//}
	//// c.timer is also immutable (it is set after make(chan) but before any channel operations).
	//// All timer channels have dataqsiz > 0.
	//if c.timer != nil {
	//	c.timer.maybeRunChan()
	//}
	//return atomic.Loaduint(&c.qcount) == 0
	return false
}