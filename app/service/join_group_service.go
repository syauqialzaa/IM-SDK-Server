package service

import (
	"gin-chat-svc/app/model"
	"gin-chat-svc/pkg/common/constant"
	"gin-chat-svc/pkg/common/response"
	"gin-chat-svc/pkg/misprint"
	"time"

	"github.com/google/uuid"
)

func (g *GroupService) SaveGroup(userUuid string, group *model.Group) (model.Group, error) {
	// db := db.GetDB()
	var fromUser *model.User

	Db.Find(&fromUser, "uuid = ?", userUuid)
	if fromUser.ID == constant.NULL_ID {
		return *group, misprint.New("user doesn't exist")
	}

	// group.CreatedAt = time.Now().Local()
	// group.UpdatedAt = timeZero
	// group.DeletedAt = timeZero

	// create unique group UUID
	group.Uuid = (constant.GROUP_TAG_UUID + uuid.New().String())
	
	group.CreatorUuid = fromUser.Uuid
	group.Creator = fromUser.Username
	
	Db.Save(&group)

	groupMember := model.GroupMember {
		JoinedAt: 	time.Now().Local(),
		GroupUuid: 	group.Uuid,
		Name: 		group.Name,
		UserUuid: 	fromUser.Uuid,
		Username: 	fromUser.Username,
		Mute: 		0,
	}

	Db.Save(&groupMember)

	// also store to all_user_interacts model
	// =====================================================================
	storeInteractWithGroup := model.AllUserInteract {
		CreatedAt: 			time.Now().Local(),
		UserUuid: 			fromUser.Uuid,
		Username: 			fromUser.Username,
		InteractWithUuid: 	group.Uuid,
		InteractWith: 		group.Name,
	}

	Db.Save(&storeInteractWithGroup)
	// =====================================================================

	// also store to user_interact_visibilty model
	groupInteractVb := model.UserInteractVisibility {
		UserUuid: 			fromUser.Uuid,
		InteractWithUuid: 	group.Uuid,
		WhosLeave: 		"",
		DeleteStatus: 		"",
	}

	Db.Save(&groupInteractVb)

	return *group, nil
}

func (g *GroupService) JoinGroup(groupUuid, userUuid string) (model.GroupMember, error) {
	// db := db.GetDB()
	var user model.User
	var group model.Group
	var groupMember model.GroupMember

	Db.First(&user, "uuid = ?", userUuid)
	if user.ID == constant.NULL_ID {
		return groupMember, misprint.New("user doesn't exist")
	}

	Db.First(&group, "uuid = ?", groupUuid)
	if user.ID == constant.NULL_ID {
		return groupMember, misprint.New("group doesn't exist")
	}

	groupMemberInsert := model.GroupMember {
		JoinedAt: 	time.Now().Local(),
		UserUuid: 	user.Uuid,
		GroupUuid: 	group.Uuid,
		Name: 		group.Name,
		Username: 	user.Username,
		Mute: 		0,
	}

	// check if user already member of a group
	Db.First(&groupMember, "user_uuid = ? AND group_uuid = ? AND left_at IS NULL", user.Uuid, group.Uuid)
	if groupMember.ID != constant.NULL_ID {
		return groupMember, misprint.New("already joined the group")
	}

	// rejoin group/reinteract group
	Db.First(&groupMember, "user_uuid = ? AND group_uuid = ? AND left_at IS NOT NULL", user.Uuid, group.Uuid)
	if groupMember.ID != constant.NULL_ID {
		Db.Model(&groupMember).Update("joined_at", time.Now().Local())
		Db.Raw(`UPDATE group_members SET left_at = NULL WHERE user_uuid = ? AND group_uuid = ?`, 
			user.Uuid, group.Uuid).Scan(&groupMember)
		// Db.Model(&groupMember).Update("deleted_at", timeZero)
	}

	// add to user_interact_visibilty model
	var interactVb *model.UserInteractVisibility
	Db.First(&interactVb, "user_uuid = ? AND interact_with_uuid = ?", user.Uuid, group.Uuid)
	if interactVb.ID != constant.NULL_ID {
		Db.Model(&interactVb).Update("delete_status", constant.REINTERACT)
	} else if interactVb.ID == constant.NULL_ID {
		groupInteractVb := model.UserInteractVisibility {
			UserUuid: 			user.Uuid,
			InteractWithUuid: 	group.Uuid,
			WhosLeave: 		"",
			DeleteStatus: 		"",
		}

		Db.Save(&groupInteractVb)
	}


	// also store to all_user_interacts model
	// =====================================================================
	var allInteract *model.AllUserInteract
	Db.First(&allInteract, "user_uuid = ? AND interact_with_uuid = ? AND deleted_at IS NULL", user.Uuid, group.Uuid)
	if allInteract.ID != constant.NULL_ID {
		return groupMember, misprint.New("already joined with this group")
	}

	// rejoin/reinteract
	Db.First(&allInteract, "user_uuid = ? AND interact_with_uuid = ? AND deleted_at IS NOT NULL", user.Uuid, group.Uuid)
	if allInteract.ID != constant.NULL_ID {
		Db.Model(&allInteract).Update("created_at", time.Now().Local())
		Db.Raw(`UPDATE all_user_interacts SET deleted_at = NULL WHERE user_uuid = ? AND interact_with_uuid = ?`,
			user.Uuid, group.Uuid).Scan(&allInteract)
		return groupMember, nil		// reinteract stop right here
	}

	storeInteractWithGroup := model.AllUserInteract {
		CreatedAt: 			time.Now().Local(),
		UserUuid: 			user.Uuid,
		Username: 			user.Username,
		InteractWithUuid: 	group.Uuid,
		InteractWith: 		group.Name,
	}

	Db.Save(&storeInteractWithGroup)
	// =====================================================================

	Db.Save(&groupMemberInsert)

	return groupMemberInsert, nil
}

func (g *GroupService) LeaveGroup(groupUuid, userUuid string) (response.GroupInfo, error) {
	var user model.User
	var groupInfo response.GroupInfo

	Db.First(&user, "uuid = ?", userUuid)
	if user.ID == constant.NULL_ID {
		return groupInfo, misprint.New("user doesn't exist")
	}

	var group model.Group
	Db.First(&group, "uuid = ?", groupUuid)
	if user.ID == constant.NULL_ID {
		return groupInfo, misprint.New("group doesn't exist")
	}

	var groupMember model.GroupMember
	Db.First(&groupMember, "user_uuid = ? AND group_uuid = ? AND left_at IS NULL", user.Uuid, group.Uuid)
	if groupMember.ID == constant.NULL_ID {
		return groupInfo, misprint.New("user not joined the group yet")
	}

	// also check to user_interact_visibilities model
	// =====================================================================
	var allInteract *model.AllUserInteract

	Db.First(&allInteract, "user_uuid = ? AND interact_with_uuid = ? AND deleted_at IS NULL", user.Uuid, group.Uuid)
	if allInteract.ID == constant.NULL_ID {
		return groupInfo, misprint.New("not interact with this group yet")
	}
	
	Db.Model(&allInteract).Update("deleted_at", time.Now().Local())
	
	// =====================================================================

	var interactVb *model.UserInteractVisibility
	Db.First(&interactVb, "user_uuid = ? AND interact_with_uuid = ?", user.Uuid, group.Uuid)
	if interactVb.ID != constant.NULL_ID {
		Db.Model(&interactVb).Update("last_left", time.Now().Local())
		Db.Model(&interactVb).Update("whos_leave", user.Username)
		Db.Model(&interactVb).Update("delete_status", constant.DELETE_FOR_SELF)
	}

	Db.Model(&groupMember).Update("left_at", time.Now().Local())
	
	members := MemberList(groupUuid)
	grpInfo := response.GroupInfo {
		ID: 		group.ID,
		Name: 		group.Name,
		Creator: 	group.Creator,
		Members: 	members,
	}

	return grpInfo, nil
}