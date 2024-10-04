package ctx

type valueCtx struct {
	Context
	key, val any
}

// WithValue 基于父Context生成一个新的Context, 保存了一个key-value对
func WithValue(parent Context, key, val any) Context {
	return &valueCtx{parent, key, val}
}

// 覆盖了Value的方法,优先从自己存储中查找key.
// 实现了链式查找的功能
// 如果Context自己没有持有这个Key, 它就在其父Context中查找
func (vc *valueCtx) Value(key any) any {
	if vc.key == key {
		return vc.val
	}
	return vc.Context.Value(key)
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
