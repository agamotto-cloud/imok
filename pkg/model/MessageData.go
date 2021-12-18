package model

type MessageData struct {
	Body        string      `json:"body" `       // 消息内容
	To          string      `json:"to"   `       // 发送目标
	Sender      string      `json:"sender"`      // 发送者
	MessageType MessageType `json:"messageType"` // 消息类型
}
