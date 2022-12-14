package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func notification(ctx *gin.Context) {
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()
	for {
		msg := `{"hasError":false,"errorId":"","errorDesc":"","data":{"ruleCnt":0,"matchList":null,"total":0}}`
		err = ws.WriteMessage(websocket.TextMessage, []byte(msg))
		time.Sleep(time.Second * 10)
		log.Println(err)
	}
}

func main() {
	r := gin.Default()
	r.GET("/ws/ntf", notification)
	r.Run(":8080")
}
