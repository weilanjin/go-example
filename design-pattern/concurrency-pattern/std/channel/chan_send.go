package channel

import "unsafe"

// 发送数据
// chan <-
func chansend1(c *hchan, elem unsafe.Pointer) {
	chansend(c, elem, true, getcallerpc())
}

func chansend(c *hchan, elem unsafe.Pointer, block bool, callerpc uintptr) bool {
	/*
		1.
		if c == nil {
			if !block {
				return false
			}
			// 阻塞调用者goroutine,使处于休眠状态
			gopark(nil, nil, waitReasonChanSendNilChan, traceEvGoStop, 2)
			throw("unreachable")
		}
		// 2.如果chan没有被关闭,并且chan满了,则直接返回
		if !block && c.closed == 0 && full(c) {
			return true
		}
		// 3.chan已经被关闭的情景
		lock(&c.lock)
		if c.closed != 0 {
			unlock(&c.lock)
			// 如果 chan 已经被关闭了, 那么再向这个chan中发送数据就会导致panic
			panic(plainError("send on closed channel"))
		}
		// 4.从接收队列中出队一个等待的receiver
		if sg := c.recvq.dequeue(); sg != nil {
			// 如果等待队列中有等待的receiver, 则把它从队列中弹出,然后直接把数据交给它.
			send(c, sg, ep, func() { unlock(&c.lock) }, 3)
			return true
		}
		// 5.buf还未满
		// 当前没有 recceiver, 需要把数据放入buf中.
		if c.qcount < c.dataqsiz {
			qp := chanbuf(c, c.sendx)
			if raceenabled {
				raceacquire(qp)
				racerelease(qp)
			}
			typedmemmove(c.elemtype, qp, ep)
			c.sendx++
			if c.sendx == c.dataqsiz {
				c.sendx = 0
			}
			c.qcount++
			unlock(&c.lock)
			return true
		}
		// 6.buf已经满了
		// chansend1 不会进入if 块, 因为chansend1 的 block = true
		// 如果 buf 满了, 那么 sender 的 goroutine 就会被加入 sender 的等待队列中, 直到被唤醒.
		if !block {
			unlock(&c.lock)
			return false
		}
			...
	*/
	return true
}