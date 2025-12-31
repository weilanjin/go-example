package tcpip

import "net"

// TCP 的编程接口
// TCP 协议包含了IP数据包的序号、重传次数等信息，可以解决丢包重传，纠正乱序，确保数据可靠传输。
// 问题:
//   音视频的传输，在网络比较差的情况下，我们往往希望丢掉一些帧，但是由于TCP协议的重传机制，
//   可能会反而加剧了网络拥塞的情况。

type TCPAddr struct {
	IP   net.IP
	Port int
	Zone string // IPv6 scope zone
}

func DialTCP(network string, laddr, raddr *TCPAddr) (*net.TCPConn, error)
func ListenTCP(network string, laddr *TCPAddr) (*net.TCPListener, error)

type TCPConn struct{}

func (c *TCPConn) Read(b []byte) (int, error)
func (c *TCPConn) Write(b []byte) (int, error)
func (c *TCPConn) Close() error

type TCPListener struct{}

func (l *TCPListener) Accept() (*TCPConn, error)
func (l *TCPListener) AcceptTCP() (*TCPConn, error)
func (l *TCPListener) Close() error
