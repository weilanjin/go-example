package iface

import "runtime/debug"

// Logger is logger interface
type Logger interface {
	Printf(string, ...any)
}

// LoggerFunc is a bridge between Logger and any third party logger
type LoggerFunc func(string, ...any)

// Printf implements Logger interface
func (f LoggerFunc) Printf(format string, args ...any) {
	f(format, args...)
}

// Recovery catch go runtime panic
func Recovery(logger Logger) {
	if err := recover(); err != nil {
		logger.Printf("handle recovery error: %v", err)
		logger.Printf("full_stack:%s", string(debug.Stack()))
	}
}