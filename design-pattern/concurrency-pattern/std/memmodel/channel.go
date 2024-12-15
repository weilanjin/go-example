package memmodel

// 规则1: 向一个channel中发送数据, 一定 synchronized before 对应着从这个channel中接收数据的完成

func sendChan() {
	var ch = make(chan struct{}, 10) // buffer channel or unbuffered channel
	var s string

	var f = func() {
		s = "hello world"
		ch <- struct{}{}
	}

	go f()
	<-ch
	println(s) // s is "hello world"
}

// 规则2: chanel 的关闭完成,一定 synchronized before 由于 channel 关闭导致 receiver 接收到零值

func closeChan() {
	var ch = make(chan struct{}, 10) // buffer channel or unbuffered channel
	var s string

	var f = func() {
		s = "hello world"
		close(ch)
	}

	go f()
	<-ch
	println(s) // s is "hello world"
}

// 规则3: 对于 unbuffered channel, 从此channel 中读取数据的调用一定 synchronized before 向此channel中发送数据的调用完成

func unbufferedChan() {
	var ch = make(chan struct{}) //unbuffered channel
	var s string

	var f = func() {
		s = "hello world"
		<-ch
	}

	go f()
	ch <- struct{}{}
	println(s) // s is "hello world"
}

// 规则4: 如果 channel 的容量是m(m>0), 那么第n个接收操作一定 synchronized before 第n+m个发送操作完成