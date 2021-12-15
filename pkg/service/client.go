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
	socket    *websocket.Conn
	reader    *io.PipeReader
	writer    *io.PipeWriter
	closeFlag int32
}

func (m *Client) Close() error {
	atomic.StoreInt32(&m.closeFlag, -1)

	err1 := m.socket.Close()
	err2 := m.reader.Close()
	err3 := m.writer.Close()

	if err1 != nil {
		return err1
	}

	if err2 != nil {
		return err2
	}

	if err3 != nil {
		return err3
	}

	return nil
}

func (m *Client) Send(data []byte, client string) error {
	m.socket.WriteMessage(websocket.BinaryMessage, data)
	return nil
}

func (m *Client) Listen() {
	go func(m *Client) {
		for {
			_, data, err := m.socket.ReadMessage()
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				return
			}

			if err != nil {
				log.Println("fail to read message")
				continue
			}
			_, err = m.writer.Write(data)
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

func (m *Client) Write() {
	go func(m *Client) {
		var buf []byte = make([]byte, 1024)
		for {
			n, _ := m.reader.Read(buf)
			log.Println(string(buf[:n]))
		}
	}(m)
}
