package main

import (
	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewDevelopment()).Named("go")
	logger.Sugar().Infof("\x1b[%dm%s\x1b[0m", 31, "weilanjin")
}
