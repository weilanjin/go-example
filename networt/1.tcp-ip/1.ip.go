package tcpip

import "net"

// IP 协议
// 问题: IP 协议本身只定义了数据的目标IP,那么这个IP地址对应的计算机收到数据后，
//      如何将数据交给对应的应用程序呢？
// 解决:
//      在 IP 之上引入了传输层协议(TCP/UDP),通过端口号来区分不同的应用程序。
//      ip:port IP + 端口号 解决了软件间的冲突问题

type IPAddr struct {
	IP   net.IP
	Zone string // IPv6 scope zone
}

func DialIP(network string, laddr, raddr *IPAddr) (*net.IPConn, error)
func ListenIP(network string, laddr *IPAddr) (*net.IPConn, error)

type IPConn struct{}

func (c *IPConn) Read(b []byte) (int, error)
func (c *IPConn) ReadFrom(b []byte) (int, *IPAddr, error)
func (c *IPConn) ReadFromIP(b []byte) (int, *IPAddr, error)
func (c *IPConn) Write(b []byte) (int, error)
func (c *IPConn) WriteTo(b []byte, addr *IPAddr) (int, error)
func (c *IPConn) WriteToIP(b []byte, addr *IPAddr) (int, error)
func (c *IPConn) Close() error
