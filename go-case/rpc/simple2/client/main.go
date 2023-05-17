package main

import (
	"fmt"
	"lovec.wlj/go-case/rpc/simple2/client_proxy"
)

func main() {
	stub := client_proxy.NewCalcServiceClient("tcp", "localhost:10000")
	var reply int
	err := stub.Add(2, &reply)
	if err != nil {
		panic(err)
	}
	fmt.Println(reply)
}
