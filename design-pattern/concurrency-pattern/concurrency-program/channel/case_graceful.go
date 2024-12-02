package channel

import (
	"os"
	"os/signal"
	"syscall"
)

// GracefulStop 优雅退出
func GracefulStop() {
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	// 等待 Ctrl + C 中断信号
	<-termChan
}