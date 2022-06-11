package middleware

// 常见中间件 https://github.com/gin-gonic/contrib

import "net/http"

// r = NewRouter()
// r.Use(logger)
// r.Use(timeout)
// r.Use(ratelimit)
// r.Add("/", helloHandler)

type middleware func(http.Handler) http.Handler

type Router struct {
	middlewareChain []middleware
	mux             map[string]http.Handler
}

func NewRouter() *Router {
	return &Router{
		mux: make(map[string]http.Handler),
	}
}

func (r *Router) Use(m ...middleware) {
	r.middlewareChain = append(r.middlewareChain, m...)
}

func (r *Router) Add(route string, h http.Handler) {
	var mergedHandler = h
	for i := len(r.middlewareChain); i >= 0; i-- {
		mergedHandler = r.middlewareChain[i](mergedHandler)
	}
	r.mux[route] = mergedHandler
}
