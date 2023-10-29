package main

import (
	"log/slog"
	"time"

	"github.com/gorhill/cronexpr"
)

func main() {
	expr, err := cronexpr.Parse("*/2 * * * * * *")
	if err != nil {
		slog.Error(err.Error())
		return
	}
	now := time.Now()
	nextTime := expr.Next(time.Now())
	time.AfterFunc(nextTime.Sub(now), func() {
		slog.Info("hello !")
	})
	time.Sleep(5 * time.Second)
}
