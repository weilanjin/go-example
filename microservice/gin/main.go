package main

import "github.com/gin-gonic/gin"

func main() {
	e := gin.Default()
	e.POST("/upload", upload)
	e.Run(":8080")
}
