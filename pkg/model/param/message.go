package param

import "github.com/agamotto-cloud/imok/pkg/model"

type MessagePram struct {
	Body        string            `json:"body" `       // 消息内容
	To          string            `json:"to"   `       // 发送目标
	MessageType model.MessageType `json:"messageType"` // 消息类型
}
