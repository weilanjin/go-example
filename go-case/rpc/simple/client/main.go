package main

import (
	"fmt"
	"net/rpc"
)

func main() {
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
