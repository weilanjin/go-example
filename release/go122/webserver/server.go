package main

import (
	"log/slog"
	"net/http"
)

type HttpServer struct {
	addr string
}

func NewHttpServer(addr string) *HttpServer {
	return &HttpServer{
		addr: addr,
	}
}

func (s *HttpServer) Run() error {
	router := http.NewServeMux()
	router.HandleFunc("GET /user/{id}", NewUserServer().UserInfo)

	// 路由组
	v1 := http.NewServeMux()
	v1.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	// 添加中间件
	middlewareChain := MiddlewareChain(
		LoggerMiddleware,
		AuthMiddleware,
	)

	server := http.Server{
		Addr:    s.addr,
		Handler: middlewareChain(v1),
	}

	slog.Debug("http server start", slog.String("addr", s.addr))
	return server.ListenAndServe()
}
