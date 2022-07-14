package service

import (
	"gin-chat-svc/app/model"
	"gin-chat-svc/pkg/common/constant"
	"gin-chat-svc/pkg/common/request"
	"gin-chat-svc/pkg/common/response"
	"gin-chat-svc/pkg/logger"
	"gin-chat-svc/pkg/misprint"
	"gin-chat-svc/pkg/protocol"
	"time"
)

func (m *MessageService) SaveMessage(message protocol.Message) model.Message {
	// db := db.GetDB()
	var fromUser model.User
	var msg model.Message

	Db.Find(&fromUser, "uuid = ?", message.From)
	if fromUser.ID == constant.NULL_ID {
		logger.Logger.Error("service", logger.Any("user doesn't exist", fromUser.ID))
		return msg
	}

	// var toUserID int32 = 0
	var targetUuid string

	if message.MessageType == constant.MESSAGE_TYPE_USER {
		var toUser model.User
		
		Db.Find(&toUser, "uuid = ?", message.To)
		if toUser.ID == constant.NULL_ID {
			return msg
		}

		targetUuid = toUser.Uuid
		// toUserID = toUser.ID
	}

	if message.MessageType == constant.MESSAGE_TYPE_GROUP {
		var group model.Group

		Db.Find(&group, "uuid = ?", message.To)
		if group.ID == constant.NULL_ID {
			return msg
		}

		targetUuid = group.Uuid
		// toUserID = group.ID
	}

	saveMessage := model.Message {
		CreatedAt: 		time.Now().Local(),
		FromUuid: 		fromUser.Uuid,
		TargetUuid: 	targetUuid,
		Content: 		message.Content,
		ContentType: 	int16(message.ContentType),
		MessageType: 	int16(message.MessageType),
		Url: 			message.Url,
	}

	Db.Save(&saveMessage)
	return saveMessage
}

// get message based on ID and from_uuid IN and target_uuid IN (the chat room)
func (m *MessageService) GetMessageById(msgReq request.MsgRequestById) (response.MessageResponse, error) {
	var msgResp response.MessageResponse
	
	if msgReq.MessageType == constant.MESSAGE_TYPE_USER {
		var queryUser *model.User
		
		Db.First(&queryUser, "uuid = ?", msgReq.Uuid)
		if queryUser.ID == constant.NULL_ID {
			return msgResp, misprint.New("user doesn't exist")
		}
		
		var interactWith *model.User
		
		Db.First(&interactWith, "uuid = ?", msgReq.InteractWith)
		if interactWith.ID == constant.NULL_ID {
			return msgResp, misprint.New("user doesn't exist")
		}

		var msg *model.Message

		Db.First(&msg, "id = ?", msgReq.ID)
		if msg.ID == constant.NULL_ID {
			return msgResp, misprint.New("message doesn't exist")
		}

		Db.Raw(`
		SELECT 
			m.id, 
			m.from_uuid,
			m.target_uuid, 
			m.content, 
			m.content_type, 
			m.url, 
			m.created_at, 
			u.username 
		AS 
			from_username, 
			u.avatar, 
			to_user.username 
		AS 
			to_username 
		FROM messages AS m LEFT JOIN users AS u ON m.from_uuid = u.uuid 
		LEFT JOIN users AS to_user ON m.target_uuid = to_user.uuid 
		WHERE m.deleted_at IS NULL AND m.id = ? AND from_uuid IN (?, ?)
		AND target_uuid IN (?, ?)
		`, msg.ID, queryUser.Uuid, interactWith.Uuid, queryUser.Uuid, interactWith.Uuid,
		).Scan(&msgResp)

		return msgResp, nil

	} else if msgReq.MessageType == constant.MESSAGE_TYPE_GROUP {
		var queryUser *model.User
		
		Db.First(&queryUser, "uuid = ?", msgReq.Uuid)
		if queryUser.ID == constant.NULL_ID {
			return msgResp, misprint.New("user doesn't exist")
		}
		
		var interactGroup *model.Group
		
		Db.First(&interactGroup, "uuid = ?", msgReq.InteractWith)
		if interactGroup.ID == constant.NULL_ID {
			return msgResp, misprint.New("group doesn't exist")
		}

		var msg *model.Message

		Db.First(&msg, "id = ?", msgReq.ID)
		if msg.ID == constant.NULL_ID {
			return msgResp, misprint.New("message doesn't exist")
		}

		Db.Raw(`
		SELECT 
			m.id, 
			m.from_uuid,
			m.target_uuid, 
			m.content, 
			m.content_type, 
			m.url, 
			m.created_at, 
			u.username 
		AS 
			from_username, 
			u.avatar, 
			to_group.name 
		AS 
			to_group_name
		FROM messages AS m LEFT JOIN users AS u ON m.from_uuid = u.uuid 
		LEFT JOIN groups AS to_group ON m.target_uuid = to_group.uuid 
		WHERE m.deleted_at IS NULL AND m.id = ? AND from_uuid IN (?, ?) 
		AND target_uuid IN (?, ?)
		`, msg.ID, queryUser.Uuid, interactGroup.Uuid, queryUser.Uuid, interactGroup.Uuid,
		).Scan(&msgResp)

		return msgResp, nil
	}

	return msgResp, misprint.New("query type not supported")
}

func (m *MessageService) ModifyMessage(message *model.Message) (model.Message, error) {
	if err := Db.Model(&message).Updates(map[string]interface{} {
		"content":		message.Content,
		"send_time":	message.SendTime,
		"sent_time":	message.SentTime,
		"read_time":	message.ReadTime,
	}).Error; err != nil {
		return *message, misprint.New("failed to update.")
	}

	Db.Model(&message).Update("updated_at", time.Now().Local())

	return *message, nil
}

// message visibility, to handle delete message
func (m *MessageService) DeleteMsgForAll(msgReq request.MsgRequestById) error {
	var fromUuid *model.User
	var targetUuid string
	
	Db.First(&fromUuid, "uuid = ?", msgReq.Uuid)
	if fromUuid.ID == constant.NULL_ID {
		return misprint.New("user doesn't exist")
	}

	var message *model.Message
	
	Db.First(&message, "id = ?", msgReq.ID)
	if message.ID == constant.NULL_ID {
		return misprint.New("message doesn't exist")
	}

	// Db.First(&message, "message_type = ?", msgReq.MessageType)
	if msgReq.MessageType == constant.MESSAGE_TYPE_USER {
		var interactUser *model.User

		Db.First(&interactUser, "uuid = ?", msgReq.InteractWith)
		if interactUser.ID == constant.NULL_ID {
			return misprint.New("user doesn't exist")
		}

		targetUuid = interactUser.Uuid
	
	} else if msgReq.MessageType == constant.MESSAGE_TYPE_GROUP {
		var interactGroup *model.Group

		Db.First(&interactGroup, "uuid = ?", msgReq.InteractWith)
		if interactGroup.ID == constant.NULL_ID {
			return misprint.New("group doesn't exist")
		}

		targetUuid = interactGroup.Uuid
	}

	if message.FromUuid == fromUuid.Uuid {
		Db.Model(&message).Update("delete_for", constant.DELETE_FOR_ALL)

		// save to message_visibilities
		msgVisibility := model.MessageVisibility {
			CreatedAt: 		time.Now().Local(),
			MessageId: 		message.ID,
			FromUuid: 		fromUuid.Uuid,
			TargetUuid: 	targetUuid,
			WhosDelete: 	fromUuid.Username,
			DeleteStatus: 	constant.DELETE_FOR_ALL,
		}
	
		Db.Save(&msgVisibility)
		// updating deleted_at
		Db.Model(&message).Update("deleted_at", time.Now().Local())

	} else {
		return misprint.New("cannot delete opposite message")
	}

	return nil
}

func (m *MessageService) DeleteMsgForSelf(msgReq request.MsgRequestById) error {
	var fromUuid *model.User
	var targetUuid string
	
	Db.First(&fromUuid, "uuid = ?", msgReq.Uuid)
	if fromUuid.ID == constant.NULL_ID {
		return misprint.New("user doesn't exist")
	}

	var message *model.Message
	
	Db.First(&message, "id = ?", msgReq.ID)
	if message.ID == constant.NULL_ID {
		return misprint.New("message doesn't exist")
	}

	// Db.First(&message, "message_type = ?", msgReq.MessageType)
	if msgReq.MessageType == constant.MESSAGE_TYPE_USER {
		var interactUser *model.User

		Db.First(&interactUser, "uuid = ?", msgReq.InteractWith)
		if interactUser.ID == constant.NULL_ID {
			return misprint.New("user doesn't exist")
		}

		targetUuid = interactUser.Uuid
	
	} else if msgReq.MessageType == constant.MESSAGE_TYPE_GROUP {
		var interactGroup *model.Group

		Db.First(&interactGroup, "uuid = ?", msgReq.InteractWith)
		if interactGroup.ID == constant.NULL_ID {
			return misprint.New("group doesn't exist")
		}

		targetUuid = interactGroup.Uuid
	}

	if message.FromUuid == fromUuid.Uuid {
		Db.Model(&message).Update("delete_for", fromUuid.Username)

		// save to message_visibilities
		msgVisibility := model.MessageVisibility {
			CreatedAt: 		time.Now().Local(),
			MessageId: 		message.ID,
			FromUuid: 		fromUuid.Uuid,
			TargetUuid: 	targetUuid,
			WhosDelete: 	fromUuid.Username,
			DeleteStatus: 	constant.DELETE_FOR_SELF,
		}
	
		Db.Save(&msgVisibility)
		// Db.Delete(&message)
		// Db.Model(&message).Update("deleted_at", timeNow)

	} else {
		return misprint.New("cannot delete opposite message")
	}

	return nil
}