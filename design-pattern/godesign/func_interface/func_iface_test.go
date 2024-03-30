package iface

import (
	"fmt"
	"log/slog"
	"testing"
)

func TestHandler(t *testing.T) {
	var logger LoggerFunc = func(s string, a ...any) {
		slog.Error(fmt.Sprintf(s, a...))
	}
	defer Recovery(logger)
	panic("coon fail")
}
