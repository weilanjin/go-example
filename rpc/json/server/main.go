package main

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type CalcService struct{}

func (s *CalcService) Add(request int, reply *int) error {
	*reply = request + 10
	return nil
}

func main() {
	// 1.注册处理逻辑 handler
	_ = rpc.RegisterName("CalcService", &CalcService{})
	// 2. 实例化一个服务 server
	listen, _ := net.Listen("tcp", ":10001")
	// 3.启动服务
	for {
		conn, _ := listen.Accept()
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
