package service

import (
	"gin-chat-svc/app/model"
	"gin-chat-svc/pkg/common/response"
	"gin-chat-svc/pkg/misprint"
	"gin-chat-svc/utility/db"

	"github.com/google/uuid"
)

type GroupService struct {}

var NewGroupService = new(GroupService)

// var ctx *gin.Context
var Db = db.GetDB()

func (g *GroupService) GetGroups(uuid string) ([]response.GroupResponse, error) {
	// migGroup := &model.Group{}
	// migGroupMem := &model.GroupMember{}
	// db.GetDB().AutoMigrate(&migGroup)
	// db.GetDB().AutoMigrate(&migGroupMem)

	var queryUser *model.User
	// db := db.GetDB()

	Db.First(&queryUser, "uuid = ?", uuid)
	if queryUser.ID <= 0 {
		return nil, misprint.New("group doesn't exist")
	}

	var groups []response.GroupResponse

	Db.Raw(`
		SELECT g.id AS group_id, g.uuid, g.created_at, g.name, g.notice 
		FROM group_members AS gm LEFT JOIN groups AS g ON gm.group_id = g.id 
		WHERE gm.user_id = ?
	`, queryUser.ID).Scan(&groups)

	return groups, nil
}

func (g *GroupService) SaveGroup(userUuid string, group model.Group) {
	// db := db.GetDB()
	var fromUser model.User

	Db.Find(&fromUser, "uuid = ?", userUuid)
	if fromUser.ID <= 0 {
		return
	}

	group.UserId = fromUser.ID
	group.Uuid = uuid.New().String()
	Db.Save(&group)

	groupMember := model.GroupMember {
		UserId: 	fromUser.ID,
		GroupId: 	group.ID,
		Nickname: 	fromUser.Username,
		Mute: 		0,
	}

	Db.Save(&groupMember)
}

func (g *GroupService) GetUserIdByGroupUuid(groupUuid string) []model.User {
	// db := db.GetDB()
	var group model.Group

	Db.First(&group, "uuid = ?", groupUuid)
	if group.ID <= 0 {
		return nil
	}

	var users []model.User
	Db.Raw(`
		SELECT u.uuid, u.avatar, u.username 
		FROM groups AS g 
		JOIN group_members AS gm ON gm.group_id = g.id 
		JOIN users AS u ON u.id = gm.user_id WHERE g.id = ?
	`, group.ID).Scan(&users)

	return users
}

func (g *GroupService) JoinGroup(groupUuid, userUuid string) error {
	// db := db.GetDB()
	var user model.User

	Db.First(&user, "uuid = ?", userUuid)
	if user.ID <= 0 {
		return misprint.New("user doesn't exist")
	}

	var group model.Group
	Db.First(&group, "uuid = ?", groupUuid)
	if user.ID <= 0 {
		return misprint.New("group doesn't exist")
	}

	var groupMember model.GroupMember
	Db.First(&groupMember, "user_id = ? and group_id = ?", user.ID, group.ID)
	if groupMember.ID > 0 {
		return misprint.New("already joined the group")
	}

	nickname := user.Nickname
	if nickname == "" {
		nickname = user.Username
	}

	groupMemberInsert := model.GroupMember {
		UserId: 	user.ID,
		GroupId: 	group.ID,
		Nickname: 	nickname,
		Mute: 		0,
	}

	Db.Save(&groupMemberInsert)

	return nil
}