package log_test

import (
	"github.com/weilanjin/go-example/log/log"
	"github.com/weilanjin/go-example/log/log/depth"
	"testing"
)

func TestLogFormat(t *testing.T) {
	log.Debug("this is a debug log")
	log.Info("this is a debug log")
	log.Error("this is a debug log")
	depth.Add()
}