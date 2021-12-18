package service

import (
	"log"
	"sync"
)

type Group struct {
	GroupID       string
	clientUserIds []string
}

var groups = sync.Map{}

func (m Group) sendGroupMessage(message interface{}) {

	for _, userId := range m.clientUserIds {
		clientObj, _ := clients.Load(userId)
		if clientObj == nil {
			log.Println("用户 {} 当前为连接", userId)
			continue
		}
		client := clientObj.(Client)
		_ = client.SendMessage(message)
	}
}
