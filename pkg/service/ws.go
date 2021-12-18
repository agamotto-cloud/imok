package service

import (
	"github.com/agamotto-cloud/imok/pkg/config"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func checkOrigin(r *http.Request) bool {
	return true
}

var clients = sync.Map{}

func Connect(c *gin.Context) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  config.C.Server.ClientReadBuf,
		WriteBufferSize: config.C.Server.ClientWriteBuf,
		CheckOrigin:     checkOrigin,
		Subprotocols:    []string{c.GetHeader("Sec-WebSocket-Protocol")},
	}
	userId := c.GetHeader("userId")
	if len(userId) == 0 {
		c.Status(http.StatusBadRequest)
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("fail to connect, ", err)
		return
	}

	client := Client{
		ClientID: userId,
		wsConn:   conn,
	}

	client.Listen()

	clients.Store(userId, client)
}
