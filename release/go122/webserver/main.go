package main

import "log/slog"

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	if err := NewHttpServer(":8080").Run(); err != nil {
		panic(err)
	}
}
