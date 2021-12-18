package service

import (
	"github.com/agamotto-cloud/imok/pkg/model"
	"github.com/agamotto-cloud/imok/pkg/model/param"
	"github.com/agamotto-cloud/imok/pkg/model/response"
	"github.com/gin-gonic/gin"
	"strings"
)

func SendMessage(c *gin.Context) {
	var messageParam param.MessagePram
	if err := c.ShouldBindJSON(&messageParam); err != nil {
		response.Result(10001, gin.H{}, err.Error(), c)
		return
	}
	if len(messageParam.Body) == 0 {
		response.FailWithCodeMessage(11001, "消息体不能为空", c)
		return
	}
	if len(messageParam.To) == 0 {
		response.FailWithCodeMessage(11002, "发送目标不能为空", c)
		return
	}
	userId := c.GetHeader("userId")
	if len(userId) == 0 {
		response.FailWithCodeMessage(11003, "发送者的用户ID未指定", c)
		return
	}
	message := model.MessageData{
		Body:        messageParam.Body,
		To:          messageParam.To,
		Sender:      userId,
		MessageType: messageParam.MessageType,
	}
	err := sendGroupMessage(message)
	if err != nil {
		response.FailWithMessage("发送失败", c)
		return
	}
	response.OkWithMessage("发送成功", c)
}

// 发送消息到group中， 如果group不存在,则创建一个
func sendGroupMessage(message model.MessageData) (err error) {
	var group Group
	groupObj, _ := groups.Load(getGroupId(message))
	if groupObj != nil {
		group, _ = groupObj.(Group)
	} else {
		group, err = createGroup(message)
		if err != nil {
			return err
		}
	}
	group.sendGroupMessage(message)
	return nil
}

func createGroup(message model.MessageData) (Group, error) {
	var group Group
	var groupId = getGroupId(message)

	if message.MessageType == model.ONE_TO_ONE {
		group := Group{
			GroupID:       groupId,
			clientUserIds: []string{message.To, message.Sender},
		}
		groups.Store(groupId, group)
	}

	return group, nil

}

func getGroupId(message model.MessageData) string {
	if message.MessageType == model.GROUP_MESSAGE {
		return message.To
	}
	if strings.Compare(message.To, message.Sender) > 0 {
		return message.To + "-" + message.Sender
	} else {
		return message.Sender + "-" + message.To
	}
}
