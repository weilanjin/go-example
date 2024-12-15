package pool

import "sync"

type Request struct {
	ServiceMethod string // rpc 请求对象
	Seq           uint64
	next          *Request // Request 链表结构, 指向下一个请求
}

type Response struct { // rpc 响应对象
	Seq   uint64
	Error string
	next  *Response // Response 链表结构, 指向下一个响应
}

type Server struct {
	serverMap sync.Map

	reqLock sync.Mutex
	freeReq *Request // Request 链表

	respLock sync.Mutex
	freeResp *Response // Response 链表
}

func (server *Server) getRequest() *Request { // 获取请求对象
	server.reqLock.Lock() // 加锁
	req := server.freeReq
	if req == nil { // 没有可重用的对象,创建一个新的对象
		req = new(Request)
	} else {
		server.freeReq = req.next // 从链表中摘下表头
		*req = Request{}          // 清空这个对象
	}
	server.reqLock.Unlock()
	return req
}

func (server *Server) freeRequest(req *Request) { // 放回这个请求对象
	server.reqLock.Lock()
	req.next = server.freeReq // 放到链表的头部
	server.freeReq = req
	server.reqLock.Unlock()
}
