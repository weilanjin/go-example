package sync

import (
	"net"
	"testing"
)

func TestOnce(t *testing.T) {
	initConnFn := OnceFn(func() net.Conn {
		conn, err := net.Dial("tcp", "127.0.0.1:8080")
		if err != nil {
			panic(err)
		}
		return conn
	})
	conn := initConnFn()
	_ = conn
}
