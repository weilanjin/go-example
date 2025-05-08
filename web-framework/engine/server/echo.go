package server

import (
	"bufio"
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/weilanjin/go-example/web-framework/engine/pkg/atomic"
	"github.com/weilanjin/go-example/web-framework/engine/pkg/logger"
	"github.com/weilanjin/go-example/web-framework/engine/pkg/wait"
)

type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean
}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

func (c *EchoClient) Close() error {
	c.Waiting.WaitWithTimeout(10 * time.Second) //等待10s强制关闭
	err := c.Conn.Close()
	return err
}

func (h *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		_ = conn.Close()
	}
	client := &EchoClient{
		Conn: conn,
	}
	h.activeConn.Store(client, struct{}{})
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF { // 已经读完
				logger.Info("connection close")
				h.activeConn.Delete(client)
			} else {
				logger.Error(err)
			}
			return
		}
		client.Waiting.Add(1)
		conn.Write([]byte(msg)) // 处理业务
		client.Waiting.Done()
	}
}

func (h *EchoHandler) Close() error {
	logger.Info("handler shutting down....")
	h.closing.Set(true)
	h.activeConn.Range(func(key, value any) bool {
		client := key.(*EchoClient)
		_ = client.Close()
		return true // 继续处理下一个
	})
	return nil
}
