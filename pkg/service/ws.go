package service

import (
	"github.com/agamotto-cloud/imok/pkg/config"
	"github.com/agamotto-cloud/imok/pkg/model"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func checkOrigin(r *http.Request) bool {
	return true
}

var clients = sync.Map{}
var once sync.Once

func init() {
	once.Do(pingClient)
}

func Connect(c *gin.Context) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  config.C.Server.ClientReadBuf,
		WriteBufferSize: config.C.Server.ClientWriteBuf,
		CheckOrigin:     checkOrigin,
		Subprotocols:    []string{c.GetHeader("Sec-WebSocket-Protocol")},
	}
	userId, _ := getUserInfo(c)
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

func getUserInfo(c *gin.Context) (userId string, userName string) {
	userId = c.GetHeader("userId")
	if len(userId) == 0 {
		userId = c.GetHeader(model.HEADER_USER_ID_KEY)
	}
	if len(userId) == 0 {
		userId = c.GetHeader(model.HEADER_ACCOUNT_ID_KEY)
	}
	userName = c.GetHeader("userName")
	if len(userId) == 0 {
		userName = c.GetHeader(model.HEADER_USER_NAME_KEY)
	}
	if len(userId) == 0 {
		userName = c.GetHeader(model.HEADER_ACCOUNT_NAME_KEY)
	}
	return userId, userName
}

// 检查所有连接的客户端，踢出不健康的客户端
func pingClient() {
	go func() {
		for {
			currentTime := time.Now()
			clients.Range(func(key, clientObj interface{}) bool {
				client := clientObj.(Client)
				client.HealthCheck()
				//todo 去除
				return true
			})
			nextTime := currentTime.Add(1 * time.Second)
			sleepTime := nextTime.Sub(time.Now()).Nanoseconds()
			if sleepTime > 100 {
				time.Sleep(time.Duration(sleepTime))
			}
		}
	}()
}
