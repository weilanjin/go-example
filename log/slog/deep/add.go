package deep

import "log/slog"

func Add(a, b int) {
	slog.Debug("a + b = ", slog.Int("value", a+b))
}