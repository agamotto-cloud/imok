package service

import (
	"errors"
	"io"
	"log"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

type Client struct {
	ClientID  string
	wsConn    *websocket.Conn
	closeFlag int32
}

func (m *Client) Close() error {
	atomic.StoreInt32(&m.closeFlag, -1)

	err1 := m.wsConn.Close()

	if err1 != nil {
		return err1
	}

	return nil
}

func (m *Client) Send(data []byte, client string) error {
	m.wsConn.WriteMessage(websocket.BinaryMessage, data)
	return nil
}
func (m *Client) SendMessage(message interface{}) error {
	m.wsConn.WriteJSON(message)
	return nil
}

// Listen 监听客户端发来的消息, 这里用于收心跳
func (m *Client) Listen() {
	go func(m *Client) {
		for {
			_, data, err := m.wsConn.ReadMessage()
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				return
			}

			if err != nil {

				log.Println("fail to read message", err)
				continue
			}
			err = m.wsConn.WriteMessage(websocket.TextMessage, data)
			if errors.Is(err, io.ErrClosedPipe) {
				return
			}

			if err != nil {
				log.Fatalln("fail to write data")
			}

			if atomic.LoadInt32(&m.closeFlag) == -1 {
				return
			}
		}
	}(m)
}

/*func (m *Client) Write() {
	go func(m *Client) {
		var buf = make([]byte, 1024)
		for {
			n, _ := m.reader.Read(buf)
			log.Println(string(buf[:n]))
		}
	}(m)
}
*/
