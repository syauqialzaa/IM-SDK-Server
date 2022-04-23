package service

import (
	"gin-chat-svc/app/model"
	"gin-chat-svc/pkg/common/constant"
	"gin-chat-svc/pkg/common/request"
	"gin-chat-svc/pkg/common/response"
	"gin-chat-svc/pkg/logger"
	"gin-chat-svc/pkg/misprint"
	"gin-chat-svc/pkg/protocol"

	"gorm.io/gorm"
)

const NULL_ID int32 = 0

type MessageService struct {}

var NewMessageService = new(MessageService)

func (m *MessageService) GetMessages(message request.MessageRequest) ([]response.MessageResponse, error) {
	// migrate := &model.Message{}
	// db.GetDB().AutoMigrate(&migrate)
	// db := db.GetDB()

	if message.MessageType == constant.MESSAGE_TYPE_USER {
		var queryUser *model.User

		Db.First(&queryUser, "uuid = ?", message.Uuid)
		if NULL_ID == queryUser.ID {
			return nil, misprint.New("user doesn't exist")
		}

		var friend *model.User

		Db.First(&friend, "username = ?", message.FriendUsername)
		if NULL_ID == friend.ID {
			return nil, misprint.New("user doesn't exist")
		}

		var messages []response.MessageResponse

		Db.Raw(`
			SELECT 
				m.id, 
				m.from_user_id, 
				m.to_user_id, 
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
			FROM messages AS m LEFT JOIN users AS u ON m.from_user_id = u.id 
			LEFT JOIN users AS to_user ON m.to_user_id = to_user.id 
			WHERE from_user_id IN (?, ?) AND to_user_id IN (?, ?)
		`, queryUser.ID, friend.ID, queryUser.ID, friend.ID).Scan(&messages)

		return messages, nil
	}

	if message.MessageType == constant.MESSAGE_TYPE_GROUP {
		messages, err := FetchGroupMessage(Db, message.Uuid)
		if err != nil {
			return nil, err
		}

		return messages, nil
	}

	return nil, misprint.New("query type not supported")
}

func FetchGroupMessage(db *gorm.DB, toUuid string) ([]response.MessageResponse, error) {
	var group model.Group

	db.First(&group, "uuid = ?", toUuid)
	if group.ID <= 0 {
		return nil, misprint.New("group doesn't exist")
	}

	var messages []response.MessageResponse
	db.Raw(`
		SELECT 
			m.id, 
			m.from_user_id, 
			m.to_user_id, 
			m.content, 
			m.content_type, 
			m.url, 
			m.created_at, 
			u.username 
		AS 
			from_username, 
			u.avatar 
		FROM messages AS m LEFT JOIN users AS u ON m.from_user_id = u.id 
		WHERE m.message_type = 2 AND m.to_user_id = ?
	`, group.ID).Scan(&messages)

	return messages, nil
}

func (m *MessageService) SaveMessage(message protocol.Message) {
	// db := db.GetDB()
	var fromUser model.User

	Db.Find(&fromUser, "uuid = ?", message.From)
	if NULL_ID == fromUser.ID {
		logger.Logger.Error("service", logger.Any("SaveMessage not find from user", fromUser.ID))
		return
	}

	var toUserID int32 = 0

	if message.MessageType == constant.MESSAGE_TYPE_USER {
		var toUser model.User
		
		Db.Find(&toUser, "uuid = ?", message.To)
		if NULL_ID == toUser.ID {
			return
		}

		toUserID = toUser.ID
	}

	if message.MessageType == constant.MESSAGE_TYPE_GROUP {
		var group model.Group

		Db.Find(&group, "uuid = ?", message.To)
		if NULL_ID == group.ID {
			return
		}

		toUserID = group.ID
	}

	saveMessage := model.Message {
		FromUserId: 	fromUser.ID,
		ToUserId: 		toUserID,
		Content: 		message.Content,
		ContentType: 	int16(message.ContentType),
		MessageType: 	int16(message.MessageType),
		Url: 			message.Url,
	}

	Db.Save(&saveMessage)
}