package main

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

func main() {
	Job("task-1")
	select {}
}

var rdsConfig = asynq.RedisClientOpt{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
}

func Job(taskName string) {
	// 定义一个每分钟执行一次的任务处理器
	taskProcessor := asynq.NewServer(rdsConfig, asynq.Config{
		Concurrency: 10, // 并发处理的任务数量
	})
	// 注册任务处理器
	mux := asynq.NewServeMux()
	mux.HandleFunc(taskName, func(ctx context.Context, t *asynq.Task) error {
		fmt.Println("Task executed at:", t.ResultWriter().TaskID(), time.Now())
		// 在这里执行你的定时任务逻辑
		return nil
	})

	// 启动任务处理器
	taskProcessor.Run(mux)
}
