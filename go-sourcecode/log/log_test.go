package log_test

import (
	"lovec.wlj/go-sourcecode/log"
	"testing"
)

func TestLogFormat(t *testing.T) {
	log.Debug("this is a debug log")
	log.Info("this is a debug log")
	log.Error("this is a debug log")
}