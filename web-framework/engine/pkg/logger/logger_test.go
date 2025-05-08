package logger_test

import (
	"testing"

	"github.com/weilanjin/go-example/web-framework/engine/pkg/logger"
)

func TestLogger(t *testing.T) {
	logger.Debug("this is debug")
}

func TestSetup(t *testing.T) {
	logger.Setup(&logger.Settings{
		Path:       "./logs",
		Name:       "server",
		Ext:        "log",
		TimeFormat: "2006-01-02 15:04:05.00",
	})
	logger.Info("this is info")
}
