package service

import (
	"gin-chat-svc/app/model"
	"gin-chat-svc/pkg/common/request"
	"gin-chat-svc/pkg/logger"
	"gin-chat-svc/pkg/misprint"
	"time"
)

// func (m *MessageService) ForwardMessage(forwardReq request.ForwardMsgReq) (*model.Message, error) {

// 	var message *model.Message
// 	Db.First(&message, "id = ? AND from_uuid = ? AND target_uuid = ? AND message_type = ?", 
// 		forwardReq.ID,
// 		forwardReq.Uuid,
// 		forwardReq.InteractWith,
// 		forwardReq.MessageType,
// 	)
// 	if message.ID == constant.NULL_ID {
// 		return nil, misprint.New("message doesn't exist")
// 	}

// 	logger.Logger.Debug("service", logger.Any("forwarding message to", forwardReq.ForwardTo))

// 	if forwardReq.Uuid == forwardReq.ForwardTo {
// 		return message, misprint.New("cannot forwarding to your self")
// 	}

// 	fwdMessage := model.Message {
// 		CreatedAt: 		time.Now().Local(),
// 		FromUuid: 		message.FromUuid,
// 		TargetUuid: 	forwardReq.ForwardTo,
// 		Content: 		message.Content,
// 		MessageType: 	message.MessageType,
// 		ContentType: 	message.ContentType,
// 		Url: 			message.Url,
// 	}

// 	Db.Save(&fwdMessage)

// 	return &fwdMessage, nil
// }

func (m *MessageService) ForwardMessage(msgReq request.MsgRequestById, targets request.Targets) (*model.Message, error) {
	
	getMsg, err := m.GetMessageById(msgReq)
	if err != nil {
		return getMsg, misprint.New("failed to retrieve message")
	}

	var message model.Message
	for _, target := range targets.Targets {
		logger.Logger.Debug("service", logger.Any("forwarding message to", target))

		message = model.Message {
			CreatedAt: 		time.Now().Local(),
			FromUuid: 		getMsg.FromUuid,
			TargetUuid: 	target,
			Content: 		getMsg.Content,
			MessageType: 	getMsg.MessageType,
			ContentType: 	getMsg.ContentType,
			Url: 			getMsg.Url,
		}

		Db.Save(&message)
	}

	return &message, nil
}