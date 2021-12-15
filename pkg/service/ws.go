package service

import (
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/agamotto-cloud/imok/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func checkOrigin(r *http.Request) bool {
	return true
}

var clients = sync.Map{}

func Regist(c *gin.Context) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  config.C.Server.ClientReadBuf,
		WriteBufferSize: config.C.Server.ClientWriteBuf,
		CheckOrigin:     checkOrigin,
		Subprotocols:    []string{c.GetHeader("Sec-WebSocket-Protocol")},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("fail to connect, ", err)
		return
	}

	r, w := io.Pipe()

	client := Client{
		ClientID: uuid.New().String(),
		socket:   conn,
		reader:   r,
		writer:   w,
	}

	client.Listen()
	client.Write()

	clients.Store(client.ClientID, client)
}
