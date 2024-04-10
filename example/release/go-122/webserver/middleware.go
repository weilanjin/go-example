package main

import (
	"log/slog"
	"net/http"
	"time"
)

type Middleware func(next http.Handler) http.HandlerFunc

// 注意中间的顺序
func MiddlewareChain(m ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		h := next
		for i := len(m) - 1; i >= 0; i-- {
			h = m[i](h)
		}
		return h.ServeHTTP
	}
}

func LoggerMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		next.ServeHTTP(w, r)
		slog.InfoContext(r.Context(), "Request received",
			slog.String("time", time.Since(now).String()),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
		)
	})
}

// curl -H "Authorization: 123" localhost:8080/user/123
func AuthMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		// check auth
		if auth == "" {
			slog.Error("AuthMiddleware", slog.String("error", "no auth"))
			http.Error(w, "no auth", http.StatusUnauthorized)
			return
		}
		slog.Debug("AuthMiddleware", slog.String("auth", auth))
		next.ServeHTTP(w, r)
	})
}
