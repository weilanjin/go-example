package ctx

type valueCtx struct {
	Context
	key, val any
}

// WithValue 基于父Context生成一个新的Context, 保存了一个key-value对
func WithValue(parent Context, key, val any) Context {
	// .....
	return &valueCtx{parent, key, val}
}

// 覆盖了Value的方法,优先从自己存储中查找key.
// 实现了链式查找的功能
// 如果Context自己没有持有这个Key, 它就在其父Context中查找
func (vc *valueCtx) Value(key any) any {
	if vc.key == key {
		return vc.val
	}
	return value(vc.Context, key)
}

// 一直往上找, 找到一层就对比一下, 不匹配就再往上找, 找到尽头
func value(c Context, key any) any {
	for {
		switch ctx := c.(type) {
		case *valueCtx:
			if key == ctx.key { // 并且就是要查找的key, 则返回此key对应的值
				return ctx.val
			}
			c = ctx.Context
		case *cancelCtx: // 如果是 cancel Context
			if key == &cancelCtxKey {
				return c
			}
			c = ctx.Context
		case *timerCtx:
			if key == &cancelCtxKey {
				return ctx.cancelCtx
			}
			c = ctx.Context
		case *emptyCtx: // 如果是 context.Background() 或者 context.TODO() 则返回 nil
			return nil
		default: // 其他情况, 比如 自定义Context, 则调用 Value 查找, Value 的逻辑自己实现
			return ctx.Value(key)
		}
	}
}

var (
	ServerContextKey    = &contextKey{"http-server"}
	LocalAddrContextKey = &contextKey{"local-addr"}
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "net/http context value " + k.name
}

/*
	// 获取服务器基本信息和本地地址信息的方式
	// ServerContextKey 服务器启动时设置
	// LocalAddrContextKey 建立连接时设置
	func (srv *Server)  Server(l net.Listener) error {
		......
		ctx := context.WithValue(baseCtx, ServerContextKey, srv)
	}

	func (c *conn) serve(ctx context.Context) {
		c.remoteAddr = c.rwc.RemoteAddr().String()
		ctx = context.WithValue(ctx, LocalAddrContextKey, c.rwc.LocalAddr().String())
		......
	}
*/
