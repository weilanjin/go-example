package project

import (
	"log"
	"net/rpc"
)

// 半异步/半同步 Half-Async Half-Sync
// 是一种处于异步后和同步操作的并发模式,它结合了两种并发模型的优点,以便在异步和同步操作之间平衡
// 两部分
// 一部分用于处理异步事件: 单独的线程池中处理,以确保同步操作不会阻塞主线程
// 一部分用于处理同步事件: 同步事件则由主线程或独立线程池中的线程处理
// 优点:
// 可以利用异步操作的高性能和高吞吐能力,同时可以利用同步操作的简单性和易用性

// net/rpc/client.go
// 请求封装成 Call 对象, 并调用client.send发送给服务端, Call对象也会被放入等待队列中.

type Client struct {
	codec   rpc.ClientCodec
	pending map[uint64]*Call
}

type Call struct {
	ServiceMethod string     // The name of the service and method to call.
	Args          any        // The argument to the function (*struct).
	Reply         any        // The reply from the function (*struct).
	Error         error      // After completion, the error status.
	Done          chan *Call // Server sends done to tell client that server has finished.
}

func NewClientWithCodec(codec rpc.ClientCodec) *Client {
	client := &Client{
		codec:   codec,
		pending: make(map[uint64]*Call),
	}
	go client.input()
	return client
}

func (client *Client) Go(serviceMethod string, args any, reply any, done chan *Call) *Call {
	call := new(Call)
	call.ServiceMethod = serviceMethod
	call.Args = args
	call.Reply = reply
	if done == nil {
		done = make(chan *Call, 10) // buffered.
	} else {
		// If caller passes done != nil, it must arrange that
		// done has enough buffer for the number of simultaneous
		// RPCs that will be using that channel. If the channel
		// is totally unbuffered, it's best not to run at all.
		if cap(done) == 0 {
			log.Panic("rpc: done channel is unbuffered")
		}
	}
	call.Done = done
	client.send(call)
	return call
}

func (client *Client) send(call *Call) {
	// client.sending.Lock()
}

func (client *Client) input() {
}