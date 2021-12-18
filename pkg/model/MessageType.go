package model

type MessageType int32

const (
	ONE_TO_ONE    MessageType = 0 //单聊
	GROUP_MESSAGE MessageType = 1 //群消息
	ALL           MessageType = 2 //全体消息
)
