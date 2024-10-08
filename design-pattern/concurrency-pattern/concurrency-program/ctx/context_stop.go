package ctx

type stopCtx struct {
	Context
	stop func() bool
}
