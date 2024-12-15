package ctx

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"
)

// curl http://localhost:8080
// curl http://localhost:8080\?token\=123456
func TestFishWebRun(t *testing.T) {
	var coonHandler = func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		ctx.Value(FishCtxKey).(*ContextInfo).Params["LocalAddrContextKey"] = r.Context().Value(http.LocalAddrContextKey)
	}

	var authHandler = func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "123456" {
			ctx.Value(FishCtxKey).(*ContextInfo).Params["Valid"] = true
		}
	}

	var b = New()

	// 添加中间件
	b.UseFunc(coonHandler)
	b.UseFunc(authHandler)

	// 添加handler
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		params := ctx.Value(FishCtxKey).(*ContextInfo).Params
		localAddr := params["LocalAddrContextKey"].(*net.TCPAddr).AddrPort()
		valid, _ := params["Valid"].(bool)
		w.Write([]byte(fmt.Sprintf("hello world. localAddr: %s, valid: %v\n", localAddr, valid)))
	})
	b.SetMux(mux)
	b.Run(":8080")
}
