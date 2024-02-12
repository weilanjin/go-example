package handler

import (
	"fmt"
	"lovec.wlj/example/err/service"
)

func GetUser() {
	err := service.QueryUser()
	fmt.Printf("%+v", err)
}