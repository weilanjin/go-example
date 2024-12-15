package ctx

import (
	"context"
	"net/http"
)

// Context 是一个保存上下文信息的对象
type ContextInfo struct {
	Params map[string]any
}

// FishCtxKey 是一个全局的 ContextKey
var FishCtxKey = contextKey{"fish_ctx"}

// NewContext 创建一个空的 Context
func NewContext() context.Context {
	return context.WithValue(context.Background(), FishCtxKey, &ContextInfo{
		Params: make(map[string]any),
	})
}

// Handler 是中间件类型, 定义了插件的方法签名
type Handler interface {
	ServerHTTP(c context.Context, w http.ResponseWriter, r *http.Request)
}

// HandlerFunc 也是一个中间件类型, 以接口的形式提供
type HandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request)

func (fn HandlerFunc) ServerHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	fn(ctx, w, r)
}

// http web 框架

type Fish struct {
	middleware []Handler
	mux        http.Handler
}

func New() *Fish {
	return &Fish{
		middleware: make([]Handler, 0),
		mux:        http.DefaultServeMux,
	}
}

// Use 添加一个中间件
func (f *Fish) Use(handler ...Handler) {
	f.middleware = append(f.middleware, handler...)
}

// UseFunc 添加一个中间件
func (f *Fish) UseFunc(handleFunc HandlerFunc) {
	f.Use(HandlerFunc(handleFunc))
}

// SetMux 设置一个 http.Handler
func (f *Fish) SetMux(mux http.Handler) {
	f.mux = mux
}

func (f *Fish) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext()
	for _, handler := range f.middleware {
		handler.ServerHTTP(ctx, w, r)
	}
	f.mux.ServeHTTP(w, r.WithContext(ctx)) // 把自定义的 Context 传递下去
}

// Run 启动一个 http 服务
func (f *Fish) Run(addr string) error {
	return http.ListenAndServe(addr, f)
}
