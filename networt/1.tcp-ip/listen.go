package tcpip

import "net"

// 服务端

func TCPServer() {
	l, err := net.Listen("tcp", "127.0.0.1:8080")
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		// go handleConnection(conn)
		_ = conn
	}
	_ = err
}

// UDP 服务端

func UDPServer() {
	addrServer := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8082}

	c, err := net.ListenUDP("udp", addrServer)
	buf := make([]byte, 1024)
	for {
		n, addr, err := c.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		// go handleUDPConnection(c, addr, buf[:n])
		_, _ = n, addr
	}
	_ = err
}
