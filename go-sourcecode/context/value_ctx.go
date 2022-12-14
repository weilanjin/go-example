package context

type valueCtx struct {
	Context
	key, val any
}