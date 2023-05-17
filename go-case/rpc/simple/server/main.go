package main

import (
	"net"
	"net/rpc"
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
	listen, _ := net.Listen("tcp", ":10000")
	// 3.启动服务
	conn, _ := listen.Accept()
	rpc.ServeConn(conn)
}
