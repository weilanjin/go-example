package main

import logger "github.com/weilanjin/go-example/log"

func main() {
	logger.Init(logger.NewOptions())
	// logger.Info("this is a test")
}