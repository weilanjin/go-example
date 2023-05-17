package main

import (
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type CalcService struct{}

func (s *CalcService) Add(request int, reply *int) error {
	*reply = request + 10
	return nil
}

func main() {
	_ = rpc.RegisterName("CalcService", &CalcService{})
	http.HandleFunc("/http-rpc", func(w http.ResponseWriter, r *http.Request) {
		var conn io.ReadWriteCloser = struct {
			io.Writer
			io.ReadCloser
		}{
			ReadCloser: r.Body,
			Writer:     w,
		}
		_ = rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
	})
	_ = http.ListenAndServe(":10002", nil)
}
