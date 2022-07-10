package service

import (
	"gin-chat-svc/app/model"
	"gin-chat-svc/pkg/common/constant"
	"gin-chat-svc/pkg/common/request"
	"gin-chat-svc/pkg/common/response"
	"gin-chat-svc/pkg/misprint"
	"sort"
	"time"

	"gorm.io/gorm"
)

func (m *MessageService) GetMessages(message request.MessageRequest) ([]response.MessageResponse, error) {
	// var queryMsg *model.Message
	
	// Db.First(&queryMsg, "message_type = ?", message.MessageType)
	if message.MessageType == constant.MESSAGE_TYPE_USER {
		var queryUser *model.User
		
		Db.First(&queryUser, "uuid = ?", message.Uuid)
		if queryUser.ID == constant.NULL_ID {
			return nil, misprint.New("user doesn't exist")
		}

		var interactWith *model.User
	
		Db.First(&interactWith, "uuid = ?", message.InteractWith)
		if interactWith.ID == constant.NULL_ID {
			return nil, misprint.New("user doesn't exist")
		}

		// get all user's message
		userMessage, err := FetchUserMessage(Db, queryUser, interactWith, queryUser.Username)
		if err != nil {
			return nil, err
		}

		// get all opposite user's message
		oppMessage, err := FetchOppMessage(Db, queryUser, interactWith)
		if err != nil {
			return nil, err
		}

		// combine user's message and opposite user's message
		listOfMessage := append(userMessage, oppMessage...)

		// sortir list of message based on message id
		sort.SliceStable(listOfMessage, func(i, j int) bool {
			return listOfMessage[i].ID < listOfMessage[j].ID
		})

		// var messages []response.MessageResponse
		
		// Db.Raw(`
		// 	SELECT 
		// 		m.id, 
		// 		m.from_uuid,
		// 		m.target_uuid, 
		// 		m.content, 
		// 		m.content_type, 
		// 		m.url, 
		// 		m.created_at, 
		// 		u.username 
		// 	AS 
		// 		from_username, 
		// 		u.avatar, 
		// 		to_user.username 
		// 	AS 
		// 		to_username 
		// 	FROM messages AS m LEFT JOIN users AS u ON m.from_uuid = u.uuid 
		// 	LEFT JOIN users AS to_user ON m.target_uuid = to_user.uuid 
		// 	WHERE m.deleted_at IS NULL AND from_uuid IN (?, ?) AND target_uuid IN (?, ?)
		// `, queryUser.Uuid, interactWith.Uuid, queryUser.Uuid, interactWith.Uuid).Scan(&messages)

		return listOfMessage, nil

	} else if message.MessageType == constant.MESSAGE_TYPE_GROUP {
		var queryUser *model.User
		
		Db.First(&queryUser, "uuid = ?", message.Uuid)
		if queryUser.ID == constant.NULL_ID {
			return nil, misprint.New("user doesn't exist")
		}
		
		var interactGroup *model.Group
		
		Db.First(&interactGroup, "uuid = ?", message.InteractWith)
		if interactGroup.ID == constant.NULL_ID {
			return nil, misprint.New("group doesn't exist")
		}
		
		// var interactWith *model.Group
	
		// Db.First(&interactWith, "name = ?", message.InteractWith)
		// if interactWith.ID == constant.NULL_ID {
		// 	return nil, misprint.New("group doesn't exist")
		// }

		userMsgInGroup, err := FetchUserMsgInGroup(Db, queryUser, interactGroup, queryUser.Username)
		if err != nil {
			return nil, misprint.New("message doesn't exist")
		}

		oppMsgInGroup, err := FetchOppMsgInGroup(Db, queryUser, interactGroup)
		if err != nil {
			return nil, misprint.New("message doesn't exist")
		}

		listOfGroupMessage := append(userMsgInGroup, oppMsgInGroup...)

		// messages, err := FetchGroupMessage(Db, message.Uuid)
		// if err != nil {
		// 	return nil, misprint.New("message doesn't exist")
		// }

		// sortir list of message based on message id
		sort.SliceStable(listOfGroupMessage, func(i, j int) bool {
			return listOfGroupMessage[i].ID < listOfGroupMessage[j].ID
		})

		return listOfGroupMessage, nil
	}

	return nil, misprint.New("query type not supported")
}

func FetchUserMessage(db *gorm.DB, queryUser, interactWith *model.User, username string) ([]response.MessageResponse, error) {
	var messages []response.MessageResponse
	var queryMsg *model.Message
	var timeParam time.Time

	var interactVb *model.UserInteractVisibility
	db.First(&interactVb, "user_uuid = ? AND interact_with_uuid = ? AND last_left IS NULL", queryUser.Uuid, interactWith.Uuid)
	if interactVb.ID != constant.NULL_ID {
		timeParam = timePast
	}

	db.First(&interactVb, "user_uuid = ? AND interact_with_uuid = ? AND last_left IS NOT NULL", queryUser.Uuid, interactWith.Uuid)
	if interactVb.ID != constant.NULL_ID {
		timeParam = interactVb.LastLeft
	}
	// if interactVb.LastLeft == nil {
		// fmt.Print(timeParam)
	// } else {
	// }

	db.First(&queryMsg, "delete_for = ?", username)
	// if delete_for is null or not queryUser username, just get all message
	if queryMsg.ID == constant.NULL_ID {
		db.Raw(`
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
		WHERE m.deleted_at IS NULL AND m.from_uuid = ? AND m.target_uuid = ?
		AND m.created_at >= ?
		`, queryUser.Uuid, interactWith.Uuid, timeParam).Scan(&messages)

		// IF(deleted_at = timeZero, timePast, userInteract.DeletedAt)FROM group_members

		// AND m.created_at >= (SELECT deleted_at,
		// IF(deleted_at IS NULL, time 2020, userInteract.DeletedAt)
		// FROM group_members)

		// AND m.created_at >= (SELECT deleted_at,
		// CASE WHEN deleted_at IS NULL THEN time 2020
		// WHEN deleted_at IS NOT NULL THEN queryMemberGroup.DeletedAt
		// END FROM group_members)
		return messages, nil
	
	} else if queryMsg.ID > constant.NULL_ID && queryMsg.DeleteFor == queryUser.Username {
		db.Raw(`
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
		WHERE m.deleted_at IS NULL AND m.delete_for NOT IN (?)
		AND m.from_uuid = ? AND m.target_uuid = ? AND m.created_at >= ?
		`, queryUser.Username, queryUser.Uuid, interactWith.Uuid, timeParam).Scan(&messages)
	}

	return messages, nil
}

// get all opposite user's message as long as deleted_at IS NULL
func FetchOppMessage(db *gorm.DB, queryUser, interactWith *model.User) ([]response.MessageResponse, error) {
	var messages []response.MessageResponse
	var interactVb *model.UserInteractVisibility
	var timeParam time.Time

	db.First(&interactVb, "user_uuid = ? AND interact_with_uuid = ? AND last_left IS NULL", queryUser.Uuid, interactWith.Uuid)
	if interactVb.ID != constant.NULL_ID {
		timeParam = timePast
	}

	db.First(&interactVb, "user_uuid = ? AND interact_with_uuid = ? AND last_left IS NOT NULL", queryUser.Uuid, interactWith.Uuid)
	if interactVb.ID != constant.NULL_ID {
		timeParam = interactVb.LastLeft
	}

	// if interactVb.LeaveAt == timeZero {
	// 	timeParam = timePast	// return TRUE, get all related message
	// } else {
	// 	timeParam = interactVb.LeaveAt
	// }

	db.Raw(`
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
	WHERE m.deleted_at IS NULL AND m.from_uuid = ?
	AND m.target_uuid = ? AND m.created_at >= ?
	`, interactWith.Uuid, queryUser.Uuid, timeParam).Scan(&messages)

	return messages, nil
}

// func FetchGroupMessage(db *gorm.DB, toUuid string) ([]response.MessageResponse, error) {
// 	var group model.Group

// 	db.First(&group, "uuid = ?", toUuid)
// 	if group.ID == constant.NULL_ID {
// 		return nil, misprint.New("group doesn't exist")
// 	}

// 	var messages []response.MessageResponse

// 	db.Raw(`
// 	SELECT 
// 		m.id, 
// 		m.from_uuid, 
// 		m.target_uuid, 
// 		m.content, 
// 		m.content_type, 
// 		m.url, 
// 		m.created_at, 
// 		u.username 
// 	AS 
// 		from_username, 
// 		u.avatar,
// 		to_group.name
// 	AS
// 		to_group_name
// 	FROM messages AS m LEFT JOIN users AS u ON m.from_uuid = u.uuid
// 	LEFT JOIN groups AS to_group ON m.target_uuid = to_group.uuid 
// 	WHERE m.message_type = 2 AND m.deleted_at IS NULL AND m.target_uuid = ?
// 	`, group.Uuid).Scan(&messages)

// 	return messages, nil
// }

func FetchUserMsgInGroup(db *gorm.DB, queryUser *model.User, interactGroup *model.Group, username string) ([]response.MessageResponse, error) {
	var messages []response.MessageResponse
	// var group model.Group
	var queryMsg *model.Message
	var timeParam time.Time
	
	var interactVb *model.UserInteractVisibility
	db.First(&interactVb, "user_uuid = ? AND interact_with_uuid = ? AND last_left IS NULL", queryUser.Uuid, interactGroup.Uuid)
	if interactVb.ID != constant.NULL_ID {
		timeParam = timePast
	}

	db.First(&interactVb, "user_uuid = ? AND interact_with_uuid = ? AND last_left IS NOT NULL", queryUser.Uuid, interactGroup.Uuid)
	if interactVb.ID != constant.NULL_ID {
		timeParam = interactVb.LastLeft
	}

	// if interactVb.LeaveAt == timeZero {
	// 	timeParam = timePast
	// } else {
	// 	timeParam = interactVb.LeaveAt
	// }

	// fmt.Println("==================================================")
	// fmt.Printf("%v\n", timeNow)

	db.First(&queryMsg, "delete_for = ?", username)
	// if delete_for is null or not queryUser username, just get all message
	if queryMsg.ID == constant.NULL_ID {
		db.Raw(`
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
		WHERE m.deleted_at IS NULL AND m.message_type = 2 AND m.from_uuid = ?
		AND m.target_uuid = ? AND m.created_at >= ?
		`, queryUser.Uuid, interactGroup.Uuid, timeParam).Scan(&messages)
		
		// AND m.created_at >= (SELECT deleted_at,
		// IF(deleted_at IS NULL, time 2020, queryMemberGroup.DeletedAt)
		// FROM group_members)

		// AND m.created_at >= (SELECT deleted_at,
		// CASE WHEN deleted_at IS NULL THEN time 2020
		// WHEN deleted_at IS NOT NULL THEN queryMemberGroup.DeletedAt
		// END FROM group_members)
		return messages, nil
	
	} else if queryMsg.ID > constant.NULL_ID && queryMsg.DeleteFor == queryUser.Username {
		db.Raw(`
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
		WHERE m.deleted_at IS NULL AND m.message_type = 2 AND m.delete_for NOT IN (?) 
		AND m.from_uuid = ? AND m.target_uuid = ? AND m.created_at >= ?
		`, queryUser.Username, queryUser.Uuid, interactGroup.Uuid, timeParam).Scan(&messages)
	}

	return messages, nil
}

// get all opposite user's message as long as deleted_at IS NULL
func FetchOppMsgInGroup(db *gorm.DB, queryUser *model.User, interactGroup *model.Group) ([]response.MessageResponse, error) {
	var messages []response.MessageResponse
	var timeParam time.Time

	var interactVb *model.UserInteractVisibility
	db.First(&interactVb, "user_uuid = ? AND interact_with_uuid = ? AND last_left IS NULL", queryUser.Uuid, interactGroup.Uuid)
	if interactVb.ID != constant.NULL_ID {
		timeParam = timePast
	}

	db.First(&interactVb, "user_uuid = ? AND interact_with_uuid = ? AND last_left IS NOT NULL", queryUser.Uuid, interactGroup.Uuid)
	if interactVb.ID != constant.NULL_ID {
		timeParam = interactVb.LastLeft
	}

	// if interactVb.LeaveAt == timeZero {
	// 	timeParam = timePast
	// } else {
	// 	timeParam = interactVb.LeaveAt
	// }

	db.Raw(`
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
	WHERE m.message_type = 2 AND m.from_uuid NOT IN (?)
	AND m.deleted_at IS NULL AND m.target_uuid = ? AND m.created_at >= ?
	`, queryUser.Uuid, interactGroup.Uuid, timeParam).Scan(&messages)
	
	return messages, nil
}