package main

import (
	"net"
	"net/rpc"
)

// CalcService
// 1.服务类型必须被导出
// 2.方法也必须被导出
// 3.方法需要正好两个参数：第一个是值类型（输入），第二个是指针（输出），这让服务器能够将结果写回给客户端
// 4.该方法只能返回一个值，这是一个错误。任何实际数据都通过第二个参数返回给客户端。
// 通信默认序列化方式：gob
// seq 的作用：
// 当客户端 建立到服务时，他们会建立一个单一的TCP连接 这个连接会一直保持打开状态，直到你明确地关闭它或者出现问题（比如服务器关闭或网络中断）
// 每个 RPC 请求都带有一个唯一的 ID 标记
// 缺点：
// 1.大量使用反射
// 2.没有使用http 2
type CalcService struct{}

func (s *CalcService) Add(request int, reply *int) error {
	*reply = request + 10
	return nil
}

func main() {

	_ = rpc.Register(new(CalcService)) // Register serviceName "CalcService"
	// 1.注册处理逻辑 handler
	// _ = rpc.RegisterName("CalcService", &CalcService{})
	// 2. 实例化一个服务 server
	listen, _ := net.Listen("tcp", ":10000")
	defer listen.Close()

	for {
		// 3.启动服务
		// 接受连接并在单独的goroutines中提供它们
		conn, _ := listen.Accept()
		// go rpc.ServeCodec(jsonrpc.NewServerCodec(conn)) // 使用json 编码
		go rpc.ServeConn(conn) // 使用 gob 编码
	}
}