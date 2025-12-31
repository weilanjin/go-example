package tcpip

import "net"

// 客户端 TCP 和 UDP 使用方式都很像
// DialTCP
// DialUDP

func main() {
	buf := make([]byte, 1024)

	addrServer := "127.0.0.1:8080"
	c, err := net.Dial("tcp", addrServer)
	c.Write(buf)
	c.Read(buf)
	c.Close()

	_ = err
}
