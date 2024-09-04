package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// type WebSocketMgr struct {
// 	conns sync.Map
// }

// func (mgr *WebSocketMgr) Store(key string, ws *websocket.Conn) {
// 	mgr.conns.Store(key, ws)
// }

// func (mgr *WebSocketMgr) Get(key string) (ws *websocket.Conn) {
// 	if conn, ok := mgr.conns.Load(key); !ok {
// 		ws_ := conn.(*websocket.Conn)
// 		// f := ws_.CloseHandler()
// 	}
// 	return
// }

// func (mgr *WebSocketMgr) Remove(key string) {
// 	mgr.conns.Delete(key)
// }

var WsConns sync.Map // 所有websocket连接集合

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func notification(ctx *gin.Context) {
	wsId := ctx.Param("id")

	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil) // http -> websocket
	if err != nil {
		log.Printf("%+v", errors.WithStack(err))
		return
	}
	defer ws.Close()

	WsConns.Store(wsId, ws)

	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			return
		}
		log.Printf("%s", p)
	}
}
