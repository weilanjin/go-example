package mutex

import (
	"log/slog"
	"testing"
)

func TestGoID(t *testing.T) {
	slog.Info("", "go id", GoID())
}
