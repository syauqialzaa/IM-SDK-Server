package ws

import (
	"gin-chat-svc/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader {
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RunWebsocket(ctx *gin.Context) {
	user := ctx.Query("user")

	if user == "" {
		return
	}

	logger.Logger.Info("websocket", zap.String("newUser", user))
	wsocket, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}

	client := &Client {
		Name: 	user,
		Conn: 	wsocket,
		Send: 	make(chan []byte),
	}

	MyServer.Register <- client
	go client.ReadPump()
	go client.WritePump()
}