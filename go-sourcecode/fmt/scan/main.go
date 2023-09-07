package main

import (
	"fmt"
	"log"
)

func main() {
	var name string
	var age int
	fmt.Println("请求输入: 姓名，年龄")
	_, err := fmt.Scanf("%s，%d", &name, &age)
	if err != nil {
		log.Println("err：", err)
		return
	}
	fmt.Println("接收23: ", name, age)
}