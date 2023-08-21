package main

import (
	"fmt"
	"lovec.wlj/web-framework/engine/config"
	"lovec.wlj/web-framework/engine/pkg/file"
	"lovec.wlj/web-framework/engine/pkg/logger"
	"lovec.wlj/web-framework/engine/server"
)

const configFile = "config.conf"

func main() {
	if file.Exists(configFile) {
		config.SetupConfig(configFile)
	} else {
		config.DefServer()
	}
	if err := server.ListenAndServeWithSignal(
		&server.Config{Address: fmt.Sprintf("%s:%d", config.ServerConfig.Bind, config.ServerConfig.Port)},
		server.NewEchoHandler(),
	); err != nil {
		logger.Fatal(err)
	}
}