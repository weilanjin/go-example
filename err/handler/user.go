package handler

import (
	"fmt"
	"github.com/weilanjin/go-example/err/service"
)

func GetUser() {
	err := service.QueryUser()
	fmt.Printf("%+v", err)
}