package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:10001")
	if err != nil {
		panic(err)
	}
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	var reply int
	err = client.Call("CalcService.Add", 2, &reply)
	if err != nil {
		panic(err)
	}
	fmt.Println(reply)
}
