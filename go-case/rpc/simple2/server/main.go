package main

import (
	"lovec.wlj/go-case/rpc/simple2/handler"
	"lovec.wlj/go-case/rpc/simple2/server_proxy"
	"net"
	"net/rpc"
)

func main() {
	// 1.注册处理逻辑 handler
	server_proxy.RegisterCalcService(&handler.CalcService{})
	// 2. 实例化一个服务 server
	listen, _ := net.Listen("tcp", ":10000")
	// 3.启动服务
	for {
		conn, _ := listen.Accept()
		go rpc.ServeConn(conn)
	}
}
