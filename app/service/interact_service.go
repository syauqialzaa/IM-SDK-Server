package service

import (
	"gin-chat-svc/app/model"
	"gin-chat-svc/pkg/common/constant"
	"gin-chat-svc/pkg/common/request"
	"gin-chat-svc/pkg/logger"
	"gin-chat-svc/pkg/misprint"
	"time"
)

func (u *UserService) AddUserInteract(userInteractRequest *request.InteractRequest) error {
	// db := db.GetDB()
	var queryUser *model.User

	Db.First(&queryUser, "uuid = ?", userInteractRequest.Uuid)
	logger.Logger.Debug("service", logger.Any("queryUser", queryUser))
	if queryUser.ID == constant.NULL_ID {
		return misprint.New("user doesn't exist")
	}

	var interactWith *model.User
	Db.First(&interactWith, "username = ?", userInteractRequest.InteractWith)	// temp, change to interactWithUuid
	if interactWith.ID == constant.NULL_ID {
		return misprint.New("interact has been added")
	}

	addUserInteract := model.UserInteract {
		CreatedAt: 			time.Now().Local(),
		// DeletedAt: 			timeZero,
		UserId: 			queryUser.ID,
		Username: 			queryUser.Username,
		InteractWithId: 	interactWith.ID,
		InteractWith: 		interactWith.Username,
	}
	
	var userInteractQuery *model.UserInteract
	Db.First(&userInteractQuery, "user_id = ? AND interact_with_id = ? AND deleted_at IS NULL", queryUser.ID, interactWith.ID)
	if userInteractQuery.ID != constant.NULL_ID {
		return misprint.New("already interact with this user")
	}

	// reinteract
	Db.First(&userInteractQuery, "user_id = ? AND interact_with_id = ? AND deleted_at IS NOT NULL", queryUser.ID, interactWith.ID)
	if userInteractQuery.ID != constant.NULL_ID {
		Db.Model(&userInteractQuery).Update("created_at", time.Now().Local())
		Db.Raw(`UPDATE user_interacts SET deleted_at = NULL WHERE user_id = ? AND interact_with_id = ?`,
			queryUser.ID, interactWith.ID).Scan(&userInteractQuery)
		// Db.Model(&userInteractQuery).Update("deleted_at", timeZero)
		// return nil
	}

	// direct add to interact vb
	var interactVb *model.UserInteractVisibility
	Db.First(&interactVb, "user_uuid = ? AND interact_with_uuid = ?", queryUser.Uuid, interactWith.Uuid)
	if interactVb.ID != constant.NULL_ID {
		Db.Model(&interactVb).Update("delete_status", constant.REINTERACT)
		// return nil
	} else if interactVb.ID == constant.NULL_ID {
		userInteractVisibility := model.UserInteractVisibility {
			UserUuid: 			queryUser.Uuid,
			InteractWithUuid: 	interactWith.Uuid,
			WhosLeave: 			"",
			DeleteStatus: 		"",
		}

		Db.Save(&userInteractVisibility)
	}

	// opp interact vb
	Db.First(&interactVb, "user_uuid = ? AND interact_with_uuid = ?", interactWith.Uuid, queryUser.Uuid)
	if interactVb.ID != constant.NULL_ID {
		return nil	// reinteract stop right here
	} else if interactVb.ID == constant.NULL_ID {
		interactOppVb := model.UserInteractVisibility {
			UserUuid: 			interactWith.Uuid,
			InteractWithUuid:	queryUser.Uuid,
			WhosLeave: 		"",
			DeleteStatus: 		"",
		}

		Db.Save(&interactOppVb)
	}

	// if has been deleted, update it
	// if userInteractQuery.DeletedAt.Before(timeNow) {
	// 	Db.Model(&userInteractQuery).Update("deleted_at", timeZero)
	// }

	Db.Save(&addUserInteract)
	logger.Logger.Debug("service", logger.Any("userInteractWith", addUserInteract))

	// also store to all_user_interacts model
	// =====================================================================
	// NewUserAndGroupService.PostInteract(userInteractRequest)
	var allInteract *model.AllUserInteract
	Db.First(&allInteract, "user_uuid = ? AND interact_with_uuid = ? AND deleted_at IS NULL", queryUser.Uuid, interactWith.Uuid)
	if allInteract.ID != constant.NULL_ID {
		return misprint.New("already interact with this user")
	}

	Db.First(&allInteract, "user_uuid = ? AND interact_with_uuid = ? AND deleted_at IS NOT NULL", queryUser.Uuid, interactWith.Uuid)
	if allInteract.ID != constant.NULL_ID {
		Db.Model(&allInteract).Update("created_at", time.Now().Local())
		Db.Raw(`UPDATE all_user_interacts SET deleted_at = NULL WHERE user_uuid = ? AND interact_with_uuid = ?`,
			queryUser.Uuid, interactWith.Uuid).Scan(&allInteract)
		// Db.Model(&allInteract).Update("created_at", time.Now().Local())
		// Db.Model(&allInteract).Update("deleted_at", timeZero)
		return nil
	}

	storeInteractWithUser := model.AllUserInteract {
		CreatedAt: 			time.Now().Local(),
		// DeletedAt: 			timeZero,
		UserUuid: 			queryUser.Uuid,
		Username: 			queryUser.Username,
		InteractWithUuid: 	interactWith.Uuid,
		InteractWith: 		interactWith.Username,
	}

	Db.Save(&storeInteractWithUser)
	logger.Logger.Debug("service", logger.Any("opposite interact", storeInteractWithUser))

	Db.First(&allInteract, "user_uuid = ? AND interact_with_uuid = ? AND deleted_at IS NULL", interactWith.Uuid, queryUser.Uuid)
	if allInteract.ID != constant.NULL_ID {
		return nil
	}

	// also add auto interact opp, so that left opp users can auto interact
	Db.First(&allInteract, "user_uuid = ? AND interact_with_uuid = ? AND deleted_at IS NOT NULL", interactWith.Uuid, queryUser.Uuid)
	if allInteract.ID != constant.NULL_ID {
		Db.Raw(`UPDATE all_user_interacts SET deleted_at = NULL WHERE user_uuid = ? AND interact_with_uuid = ?`,
			queryUser.Uuid, interactWith.Uuid).Scan(&allInteract)
		// Db.Model(&allInteract).Update("deleted_at", timeZero)
		return nil
	}

	storeOppInteract := model.AllUserInteract {
		CreatedAt: 			time.Now().Local(),
		UserUuid: 			interactWith.Uuid,
		Username: 			interactWith.Username,
		InteractWithUuid: 	queryUser.Uuid,
		InteractWith: 		queryUser.Username,
	}

	Db.Save(&storeOppInteract)
	logger.Logger.Debug("service", logger.Any("opposite interact", storeOppInteract))
	// =====================================================================

	// add user opposite interact, when user interact has added,
	// also insert the opposite interact
	addUserOppInteract := model.UserInteract {
		CreatedAt: 			time.Now().Local(),
		UserId: 			interactWith.ID,
		Username: 			interactWith.Username,
		InteractWithId: 	queryUser.ID,
		InteractWith: 		queryUser.Username,
	}

	var userOppInteractQuery *model.UserInteract
	
	Db.First(&userOppInteractQuery, "user_id = ? AND interact_with_id = ? AND deleted_at IS NULL", interactWith.ID, queryUser.ID)
	if userOppInteractQuery.ID != constant.NULL_ID {
		return nil
		// return misprint.New("user is already interact with")
	}

	Db.First(&userInteractQuery, "user_id = ? AND interact_with_id = ? AND deleted_at IS NOT NULL", interactWith.ID, queryUser.ID)
	if userInteractQuery.ID != constant.NULL_ID {
		Db.Raw(`UPDATE all_user_interacts SET deleted_at = NULL WHERE user_id = ? AND interact_with_id = ?`,
			queryUser.ID, interactWith.ID).Scan(&userInteractQuery)
		// Db.Model(&userInteractQuery).Update("deleted_at", timeZero)
		return nil
	}

	Db.Save(&addUserOppInteract)
	logger.Logger.Debug("service", logger.Any("userInteractWith", addUserOppInteract))
	
	return nil
}

func (u *UserService) DeleteUserInteract(userInteractRequest *request.InteractRequest) error {
	var queryUser *model.User

	Db.First(&queryUser, "uuid = ?", userInteractRequest.Uuid)
	if queryUser.ID == constant.NULL_ID {
		return misprint.New("user doesn't exist")
	}
	
	var interactWith *model.User
	Db.First(&interactWith, "username = ?", userInteractRequest.InteractWith)	// temp, change to interactWithUuid
	if interactWith.ID == constant.NULL_ID {
		return misprint.New("interact has been added")
	}

	var userInteractQuery *model.UserInteract
	Db.First(&userInteractQuery, "user_id = ? and interact_with_id = ? AND deleted_at IS NULL", queryUser.ID, interactWith.ID)
	if userInteractQuery.ID == constant.NULL_ID {
		return misprint.New("interact doesn't exist")
	}

	// also delete the record in user_all_interacts model
	// =====================================================================
	var allInteract *model.AllUserInteract

	Db.First(&allInteract, "user_uuid = ? AND interact_with_uuid = ? AND deleted_at IS NULL", queryUser.Uuid, interactWith.Uuid)
	if allInteract.ID == constant.NULL_ID {
		return misprint.New("interact doesn't exist")
	}

	Db.Model(&allInteract).Update("deleted_at", time.Now().Local())
	// =====================================================================
	
	// also delete the record in user_interact_visibilities model
	var interactVb *model.UserInteractVisibility
	Db.First(&interactVb, "user_uuid = ? AND interact_with_uuid = ?", queryUser.Uuid, interactWith.Uuid)
	if interactVb.ID != constant.NULL_ID {
		Db.Model(&interactVb).Update("last_left", time.Now().Local())
		Db.Model(&interactVb).Update("whos_leave", queryUser.Username)
		Db.Model(&interactVb).Update("delete_status", constant.DELETE_FOR_SELF)
	}

	Db.Model(&userInteractQuery).Update("deleted_at", time.Now().Local())
	// Db.Delete(&userInteractQuery)

	return nil
}