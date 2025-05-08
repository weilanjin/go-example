package main

import (
	"fmt"

	"github.com/weilanjin/go-example/web-framework/engine/config"
	"github.com/weilanjin/go-example/web-framework/engine/pkg/file"
	"github.com/weilanjin/go-example/web-framework/engine/pkg/logger"
	"github.com/weilanjin/go-example/web-framework/engine/server"
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
