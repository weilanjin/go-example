package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	param := map[string]any{
		"id":     0,
		"params": 2,
		"method": "CalcService.Add",
	}
	p, _ := json.Marshal(param)
	reader := bytes.NewReader(p)
	resp, err := http.Post("http://localhost:10002/http-rpc", "application/json", reader)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
