package context

import (
	"time"
)

// Context implementations
// cancelCtx
// emptyCtx
// timerCtx
// valueCtx
//
// signalCtx
// onlyValuesCtx
type Context interface {
	Deadline() (deadline time.Time, ok bool)

	// Done
	// 	func Stream(ctx context.Context, out chan<- Value) error {
	//  	for {
	//  		v, err := DoSomething(ctx)
	//  		if err != nil {
	//  			return err
	//  		}
	//  		select {
	//  		case <-ctx.Done():
	//  			return ctx.Err()
	//  		case out <- v:
	//  		}
	//  	}
	//  }
	//
	// See https://blog.golang.org/pipelines for more examples of how to use
	Done() <-chan struct{}

	// Err
	// if Done is not closed, Err returns nil
	// if Done is closed, non-nil error
	//   1. Canceled (手动取消) or DeadlineExceeded (到时间取消)
	//   2. 连续调用会返回相同的错误
	Err() error

	// Value
	// Set context.WithValue(ctx, userKey, u)
	// Get u, ok := ctx.Value(userKey).(*User)
	Value(key any) any
}