package service

import (
	"gin-chat-svc/app/model"
	"gin-chat-svc/pkg/common/constant"
	"gin-chat-svc/pkg/common/response"
	"gin-chat-svc/pkg/misprint"
)

func (g *GroupService) GetGroups(uuid string) ([]response.GroupResponse, error) {
	var queryUser *model.User

	Db.First(&queryUser, "uuid = ?", uuid)
	if queryUser.ID == constant.NULL_ID {
		return nil, misprint.New("user doesn't exist")
	}

	var groups []response.GroupResponse
	
	Db.Raw(`
		SELECT g.id AS group_id, g.uuid, g.created_at, g.name, g.notice 
		FROM group_members AS gm LEFT JOIN groups AS g ON gm.group_uuid = g.uuid 
		WHERE gm.user_uuid = ? AND gm.left_at IS NULL
	`, queryUser.Uuid).Scan(&groups)

	return groups, nil
}

func (g *GroupService) GetGroupInfo(groupUuid string) (response.GroupInfo, error) {
	var group model.Group
	var groupInfo response.GroupInfo

	Db.First(&group, "uuid = ?", groupUuid)
	if group.ID == constant.NULL_ID {
		return groupInfo, misprint.New("group doesn't exist")
	}

	members := MemberList(groupUuid)
	grpInfo := response.GroupInfo {
		ID: 		group.ID,
		Name: 		group.Name,
		Creator: 	group.Creator,
		CreatedAt: 	group.CreatedAt,
		Members: 	members,
	}

	return grpInfo, nil
}

func MemberList(groupUuid string) []response.MemberList {
	var group model.Group

	Db.First(&group, "uuid = ?", groupUuid)
	if group.ID == constant.NULL_ID {
		return nil
	}

	var memberList []response.MemberList

	Db.Raw(`
		SELECT u.uuid, u.avatar, u.username, u.nickname, u.email 
		FROM groups AS g 
		JOIN group_members AS gm ON gm.group_uuid = g.uuid 
		JOIN users AS u ON u.uuid = gm.user_uuid WHERE g.uuid = ?
		AND gm.left_at IS NULL
		`, group.Uuid).Scan(&memberList)
		
	return memberList
}

func (g *GroupService) GetUserIdByGroupUuid(groupUuid string) []model.User {
	var group model.Group
	
	Db.First(&group, "uuid = ?", groupUuid)
	if group.ID == constant.NULL_ID {
		return nil
	}

	var users []model.User

	Db.Raw(`
		SELECT u.uuid, u.avatar, u.username 
		FROM groups AS g 
		JOIN group_members AS gm ON gm.group_uuid = g.uuid 
		JOIN users AS u ON u.uuid = gm.user_uuid
		WHERE g.uuid = ? AND gm.left_at IS NULL
	`, group.Uuid).Scan(&users)

	return users
}