package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	// json codec
	// conn, err := net.Dial("tcp", "localhost:10000")
	// client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	client, err := rpc.Dial("tcp", "localhost:10000")
	if err != nil {
		panic(err)
	}
	var reply int
	err = client.Call("CalcService.Add", 2, &reply)
	if err != nil {
		panic(err)
	}
	fmt.Println(reply)
}