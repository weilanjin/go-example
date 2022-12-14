package context

// canceler implementations
// *cancelCtx
// *timerCtx
type canceler interface {
	cancel(removeFormParent bool, err error)
	Done() <-chan struct{}
}