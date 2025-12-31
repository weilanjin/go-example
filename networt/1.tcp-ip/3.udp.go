package tcpip

import "net"

// UDP 的编程接口
// 在 IP协议基础上简单的引入了端口，额外并没有做什么

type UDPAddr struct {
	IP   net.IP
	Port int
	Zone string // IPv6 scope zone
}

func DialUDP(network string, laddr, raddr *UDPAddr) (*net.UDPConn, error)
func ListenUDP(network string, laddr *UDPAddr) (*net.UDPConn, error)

type UDPConn struct{}

func (c *UDPConn) Read(b []byte) (int, error)
func (c *UDPConn) ReadFrom(b []byte) (int, *UDPAddr, error)
func (c *UDPConn) ReadFromUDP(b []byte) (int, *UDPAddr, error)
func (c *UDPConn) Write(b []byte) (int, error)
func (c *UDPConn) WriteTo(b []byte, addr *UDPAddr) (int, error)
func (c *UDPConn) WriteToUDP(b []byte, addr *UDPAddr) (int, error)
func (c *UDPConn) Close() error
