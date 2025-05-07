package trace

import (
	"context"
)

type Connector interface {
	Init(ctx context.Context) error
	Push(ctx context.Context, data ...*TraceLog) error
	Enable() bool
	Logger(err error)
}
