package main

import (
	"context"
	"log/slog"
	"os/exec"
	"time"
)

type Result struct {
	Content []byte
	Err     error
}

func main() {
	var cmd *exec.Cmd
	ch := make(chan Result, 0)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// 生成 cmd
		// cmd := exec.Command("/bin/bash", "-c", "sleep 5;ls -lh")
		cmd = exec.CommandContext(ctx, "/bin/bash", "-c", "sleep 4;ls -lh")
		// 执行命令, 捕获子进程的输出(pipe)
		res, err := cmd.Output()
		ch <- Result{Content: res, Err: err}
	}()

	time.Sleep(2 * time.Second)
	// kill pid 进程ID, 杀死子进程
	cancel()
	res := <-ch
	if res.Err != nil {
		slog.Error("cmd.Output", slog.String("sh", cmd.String() /* 命令内容 */), slog.Any("res", res))
		return
	}
	// 打印子进程输出
	slog.Info(string(res.Content))
}
