package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/weilanjin/go-example/web-framework/engine/iface/tcp"
	"github.com/weilanjin/go-example/web-framework/engine/pkg/logger"
)

type Config struct {
	Address string
}

func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	closeCh := make(chan struct{})
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-signalCh
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeCh <- struct{}{}
		}
	}()

	listen, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("bind: %s start listening...", cfg.Address))
	ListenAndServe(listen, handler, closeCh)
	return nil
}

func ListenAndServe(listener net.Listener, handler tcp.Handler, closeCh <-chan struct{}) {
	go func() {
		<-closeCh
		logger.Info("shutting down...")
		_ = listener.Close()
		_ = handler.Close()
	}()

	defer func() {
		// close during unexpected error
		_ = listener.Close()
		_ = handler.Close()
	}()
	ctx := context.Background()
	var wg sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("error accepting", err)
			break
		}
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
			handler.Handle(ctx, conn)
		}()
	}
	wg.Wait() // 等所有的请求处理完成才退出
}
