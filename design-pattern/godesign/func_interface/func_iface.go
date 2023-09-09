package iface

import "net/http"

type Handler interface {
	ServerHttp(http.ResponseWriter, *http.Request)
}

type HandlerFunc func(http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServerHttp(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}