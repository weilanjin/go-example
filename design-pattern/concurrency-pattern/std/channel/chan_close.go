package channel

// close(chan) 会编译成 closechan(chan)

// 1 如果chan = nil,则关闭它会导致panic
// 2 如果chan已关闭,则再次关闭会panic
// 3 如果chan != nil, chan 也没有关闭, 那么就把等待队列中的sender(write)和
// receiver(reader)从队列中全部移除并唤醒
func closechan(c *hchan) {
	if c == nil { // chan 为空导致panic
		panic(plainError("close of nil channel"))
	}
	lock(&c.lock)
	if c.closed != 0 { // chan 已经被关闭,导致panic
		unlock(&c.lock)
		panic(plainError("close of closed channel"))
	}
	c.closed = 1
	var glist gList

	// 释放四铺皮的reader
	for {
		sg := c.recvq.dequeue()
		// ...
		gp := sg.g
		// ...
		glist.push(gp)
	}
	// 释放所有的writer(它们会导致panic)
	for {
		sg := c.sendq.dequeue()
		// ....
		gp := sg.g
		// ...
		glist.push(gp)
	}
	unlock(&c.lock)
	for !glist.empty() {
		gp := glist.pop()
		gp.schedlink = 0
		// goready(gp, 3)
	}
}

type plainError string