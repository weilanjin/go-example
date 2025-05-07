package main

import (
	"log"

	"github.com/hibiken/asynq"
)

func main() {
	JobScheduler()
	select {}
}

var rdsConfig = asynq.RedisClientOpt{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
}

func JobScheduler() {
	scheduler := asynq.NewScheduler(rdsConfig, &asynq.SchedulerOpts{})
	_, err := scheduler.Register("@every 3s", asynq.NewTask("task-1", nil))
	if err != nil {
		log.Fatalf("Failed to register task: %v", err)
	}

	// 启动调度器
	err = scheduler.Start()
	if err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}
}
