package service

import (
	"gin-chat-svc/app/model"
	"gin-chat-svc/pkg/common/constant"
	"gin-chat-svc/pkg/common/request"
	"gin-chat-svc/pkg/common/response"
	"gin-chat-svc/pkg/logger"
	"gin-chat-svc/pkg/misprint"
	"sort"
	"strings"
	"time"
)

// NOT USED YET
func (ug *UserAndGroupService) PostInteract(interactRequest *request.AllInteractReq) (model.AllUserInteract, error) {
	var allInteract *model.AllUserInteract
	var queryUser *model.User

	Db.First(&queryUser, "uuid = ?", interactRequest.Uuid)
	if queryUser.ID == constant.NULL_ID {
		return *allInteract, misprint.New("user doesn't exist")
	}

	// post interact with user
	if strings.Contains(interactRequest.InteractWithUuid, constant.USER_TAG_UUID) {
		var interactWithUser *model.User

		Db.First(&interactWithUser, "uuid = ?", interactRequest.InteractWithUuid)
		if interactWithUser.ID == constant.NULL_ID {
			return *allInteract, misprint.New("user doesn't exist")
		}

		Db.First(&allInteract, "user_uuid = ? AND interact_with_uuid = ? AND deleted_at IS NULL", queryUser.Uuid, interactWithUser.Uuid)
		if allInteract.ID != constant.NULL_ID {
			return *allInteract, misprint.New("already interact with this user")
		}

		// if allInteract.DeletedAt.Before(timeNow) {
		// 	Db.Model(&allInteract).Update("deleted_at", timeZero)
		// }

		storeInteractWithUser := model.AllUserInteract {
			CreatedAt: 			time.Now().Local(),
			// DeletedAt: 			timeZero,
			UserUuid: 			queryUser.Uuid,
			Username: 			queryUser.Username,
			InteractWithUuid: 	interactWithUser.Uuid,
			InteractWith: 		interactWithUser.Username,
		}

		Db.Save(&storeInteractWithUser)
		logger.Logger.Debug("service", logger.Any("interact with user", storeInteractWithUser))

		// add user opposite interact, when user interact has added,
		// also insert the opposite interact
		// checking first.
		Db.First(&allInteract, "user_uuid = ? AND interact_with_uuid = ? AND deleted_at IS NULL", interactWithUser.Uuid, queryUser.Uuid)
		if allInteract.ID != constant.NULL_ID {
			return *allInteract, nil
		}

		// if allInteract.DeletedAt.Before(timeNow) {
		// 	Db.Model(&allInteract).Update("deleted_at", timeZero)
		// }

		storeOppInteract := model.AllUserInteract {
			CreatedAt: 			time.Now().Local(),
			// DeletedAt: 			timeZero,
			UserUuid: 			interactWithUser.Uuid,
			Username: 			interactWithUser.Username,
			InteractWithUuid: 	queryUser.Uuid,
			InteractWith: 		queryUser.Username,
		}

		Db.Save(&storeOppInteract)
		logger.Logger.Debug("service", logger.Any("opposite interact", storeOppInteract))
		
		return storeInteractWithUser, nil

	// post interact with group
	} else if strings.Contains(interactRequest.InteractWithUuid, constant.GROUP_TAG_UUID) {
		var interactWithGroup *model.Group

		Db.First(&interactWithGroup, "uuid = ?", interactRequest.InteractWithUuid)
		if interactWithGroup.ID == constant.NULL_ID {
			return *allInteract, misprint.New("group doesn't exist")
		}

		Db.First(&allInteract, "user_uuid = ? AND interact_with_uuid = ? AND deleted_at IS NULL", queryUser.Uuid, interactWithGroup.Uuid)
		if allInteract.ID != constant.NULL_ID {
			return *allInteract, misprint.New("already interact with this group")
		}

		// if allInteract.DeletedAt.Before(timeNow) {
		// 	Db.Model(&allInteract).Update("deleted_at", timeZero)
		// }

		storeInteractWithGroup := model.AllUserInteract {
			CreatedAt: 			time.Now().Local(),
			// DeletedAt: 			timeZero,
			UserUuid: 			queryUser.Uuid,
			Username: 			queryUser.Username,
			InteractWithUuid: 	interactWithGroup.Uuid,
			InteractWith: 		interactWithGroup.Name,
		}

		Db.Save(&storeInteractWithGroup)
		logger.Logger.Debug("service", logger.Any("opposite interact", storeInteractWithGroup))

		return storeInteractWithGroup, nil
	}

	return *allInteract, nil
}

// NOT USED YET
func (ug *UserAndGroupService) DeleteInteract(interactRequest *request.AllInteractReq) (model.AllUserInteract, error) {
	var allInteract *model.AllUserInteract
	var queryUser *model.User
	var interactUuid string

	Db.First(&queryUser, "uuid = ?", interactRequest.Uuid)
	if queryUser.ID == constant.NULL_ID {
		return *allInteract, misprint.New("user doesn't exist")
	}

	Db.First(&allInteract, "interact_with_uuid = ? AND deleted_at IS NULL", interactRequest.InteractWithUuid)
	if strings.Contains(interactRequest.InteractWithUuid, constant.USER_TAG_UUID) {
		var interactWithUser *model.User

		Db.First(&interactWithUser, "uuid = ?", interactRequest.InteractWithUuid)
		if interactWithUser.ID == constant.NULL_ID {
			return *allInteract, misprint.New("user doesn't exist")
		}
		
		Db.First(&allInteract, "user_uuid = ? AND interact_with_uuid = ? AND deleted_at IS NULL", queryUser.Uuid, interactWithUser.Uuid)
		if allInteract.ID == constant.NULL_ID {
			return *allInteract, misprint.New("not interact with this user yet")
		}

		interactUuid = interactWithUser.Uuid

	} else if strings.Contains(interactRequest.InteractWithUuid, constant.GROUP_TAG_UUID) {
		var interactWithGroup *model.Group

		Db.First(&interactWithGroup, "uuid = ?", interactRequest.InteractWithUuid)
		if interactWithGroup.ID == constant.NULL_ID {
			return *allInteract, misprint.New("group doesn't exist")
		}

		Db.First(&allInteract, "user_uuid = ? AND interact_with_uuid = ? AND deleted_at IS NULL", queryUser.Uuid, interactWithGroup.Uuid)
		if allInteract.ID == constant.NULL_ID {
			return *allInteract, misprint.New("not interact with this group yet")
		}

		interactUuid = interactWithGroup.Uuid
	}

	if allInteract.UserUuid == queryUser.Uuid {
		interactVisibility := model.UserInteractVisibility {
			LastLeft: 			time.Now().Local(),
			// UserInteractId: 	allInteract.ID,
			UserUuid: 			queryUser.Uuid,
			InteractWithUuid: 	interactUuid,
			WhosLeave:			queryUser.Username,
			DeleteStatus: 		constant.DELETE_FOR_SELF,
		}

		Db.Save(&interactVisibility)
		// Db.Delete(&allInteract)
		Db.Model(&allInteract).Update("deleted_at", time.Now().Local())
	
	} else {
		return *allInteract, misprint.New("can only delete your own interact")
	}

	return *allInteract, nil
}

// find groups or users by name
func (ug *UserAndGroupService) GetUserOrGroupByName(name string) (response.SearchResponse, error) {
	var queryUser *model.User

	Db.Select("uuid", "username", "nickname", "avatar").
		First(&queryUser, "username = ?", name)

	var queryGroup *model.Group

	Db.Select("uuid", "name").First(&queryGroup, "name = ?", name)

	search := response.SearchResponse {
		User: 	*queryUser,
		Group: 	*queryGroup,
	}

	return search, nil
}

// NOT USED YET
func (ug *UserAndGroupService) GetInteracts(uuid string) ([]response.RespUserAndGroupList, error) {
	interactWithUser, err := GetInteractWithUser(uuid)
	if err != nil {
		return nil, misprint.New("interact doesn't exist")
	}

	interactWithGroup, err := GetInteractWithGroup(uuid)
	if err != nil {
		return nil, misprint.New("interact doesn't exist")
	}

	interacts := append(interactWithUser, interactWithGroup...)
	sort.SliceStable(interacts, func(i, j int) bool {
		return interacts[i].ID < interacts[j].ID
	})

	return interacts, nil
}

// NOT USED YET
func GetInteractWithUser(uuid string) ([]response.RespUserAndGroupList, error) {
	var queryUser *model.User

	Db.First(&queryUser, "uuid = ?", uuid)
	if queryUser.ID == constant.NULL_ID {
		return nil, misprint.New("interact doesn't exist")
	}

	var interactWithUser []response.RespUserAndGroupList

	Db.Raw(`
		SELECT u.id, u.uuid, u.username, u.avatar 
		FROM all_user_interacts AS aui 
		JOIN users AS u ON aui.interact_with_uuid = u.uuid 
		WHERE aui.deleted_at IS NULL AND aui.user_uuid = ?
	`, queryUser.Uuid).Scan(&interactWithUser)
	
	return interactWithUser, nil
}

// NOT USED YET
func GetInteractWithGroup(uuid string) ([]response.RespUserAndGroupList, error) {
	var queryUser *model.User

	Db.First(&queryUser, "uuid = ?", uuid)
	if queryUser.ID == constant.NULL_ID {
		return nil, misprint.New("interact doesn't exist")
	}

	var interactWithGroup []response.RespUserAndGroupList
	
	Db.Raw(`
		SELECT g.id, g.uuid, g.name
		FROM all_user_interacts AS aui
		JOIN groups AS g ON aui.interact_with_uuid = g.uuid
		WHERE aui.deleted_at IS NULL AND aui.user_uuid = ?
	`, queryUser.Uuid).Scan(&interactWithGroup)

	return interactWithGroup, nil
}