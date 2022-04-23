package service

import (
	"gin-chat-svc/app/model"
	"gin-chat-svc/pkg/common/request"
	"gin-chat-svc/pkg/common/response"
	"gin-chat-svc/pkg/logger"
	"gin-chat-svc/pkg/misprint"
	"time"

	"github.com/google/uuid"
)

type UserService struct {}

var NewUserService = new(UserService)

func (u *UserService) Register(user *model.User) error {
	// db := db.GetDB()
	var userCount int64

	Db.Model(user).Where("username", user.Username).Count(&userCount)
	if userCount > 0 {
		return misprint.New("user already exist")
	}

	user.Uuid = uuid.New().String()
	user.CreateAt = time.Now()
	user.DeleteAt = 0

	Db.Create(&user)

	return nil
}

func (u *UserService) Login(user *model.User) bool {
	// db.GetDB().AutoMigrate(&user)
	// db := db.GetDB()

	var queryUser *model.User
	logger.Logger.Debug("service", logger.Any("user in service", user))

	Db.First(&queryUser, "username = ?", user.Username)
	logger.Logger.Debug("service", logger.Any("queryUser", queryUser))

	user.Uuid = queryUser.Uuid
	compare := queryUser.Password == user.Password
	
	return compare
}

func (u *UserService) ModifyUserInfo(user *model.User) error {
	// db := db.GetDB()
	var queryUser *model.User

	Db.First(&queryUser, "username = ?", user.Username) // user.Username
	logger.Logger.Debug("service", logger.Any("queryUser", queryUser))

	var ID_NIL int32 = 0
	if ID_NIL == queryUser.ID {
		return misprint.New("user doesn't exist")
	}

	queryUser.Nickname	= user.Nickname
	queryUser.Email		= user.Email
	queryUser.Password	= user.Password

	Db.Save(queryUser)

	return nil
}

func (u *UserService) GetUserDetails(uuid string) model.User {
	// db := db.GetDB()
	var queryUser *model.User

	Db.Select("uuid", "username", "nickname", "avatar").
		First(&queryUser, "uuid = ?", uuid)

	return *queryUser
}

// find groups or users by name
func (u *UserService) GetUserOrGroupByName(name string) response.SearchResponse {
	// db := db.GetDB()
	var queryUser *model.User

	Db.Select("uuid", "username", "nickname", "avatar").
		First(&queryUser, "username = ?", name)

	var queryGroup *model.Group
	Db.Select("uuid", "name").First(&queryGroup, "name = ?", name)

	search := response.SearchResponse {
		User: 	*queryUser,
		Group: 	*queryGroup,
	}

	return search
}

func (u *UserService) GetUserList(uuid string) []model.User {
	// db := db.GetDB()
	var queryUser *model.User

	Db.First(&queryUser, "uuid = ?", uuid)

	var ID_NIL int32 = 0
	if ID_NIL == queryUser.ID {
		return nil
	}

	var queryUsers []model.User
	Db.Raw(`
		SELECT u.username, u.uuid, u.avatar 
		FROM user_friends AS uf JOIN users AS u ON uf.friend_id = u.id 
		WHERE uf.user_id = ?
	`, queryUser.ID).Scan(&queryUsers)

	return queryUsers
}

func (u *UserService) AddFriend(userFriendRequest *request.FriendRequest) error {
	// db := db.GetDB()
	var queryUser *model.User

	Db.First(&queryUser, "uuid = ?", userFriendRequest.Uuid)
	logger.Logger.Debug("service", logger.Any("queryUser", queryUser))

	var ID_NIL int32 = 0
	if ID_NIL == queryUser.ID {
		return misprint.New("user doesn't exist")
	}

	var friend *model.User
	Db.First(&friend, "username = ?", userFriendRequest.FriendUsername)
	if ID_NIL == friend.ID {
		return misprint.New("friend has been added")
	}

	userFriend := model.UserFriend {
		UserId: 	queryUser.ID,
		FriendId: 	friend.ID,
	}

	var userFriendQuery *model.UserFriend
	Db.First(&userFriendQuery, "user_id = ? and friend_id = ?", queryUser.ID, friend.ID)
	if userFriendQuery.ID != NULL_ID {
		return misprint.New("user is already a friend")
	}

	// db.AutoMigrate(&userFriend)
	Db.Save(&userFriend)
	logger.Logger.Debug("service", logger.Any("userFriend", userFriend))

	return nil
}

// ========================================================================
func (u *UserService) DeleteFriend(id string) error {
	var userFriend *model.UserFriend

	Db.First(&userFriend, "id = ?", id)
	if NULL_ID == userFriend.ID {
		return misprint.New("user doesn't exist")
	}

	Db.Delete(&userFriend)

	return nil
}

// modify avatar
func (u *UserService) ModifyUserAvatar(avatar string, userUuid string) error {
	// db := db.GetDB()
	var queryUser *model.User

	Db.First(&queryUser, "uuid = ?", userUuid)

	if NULL_ID == queryUser.ID {
		return misprint.New("user doesn't exist")
	}

	Db.Model(&queryUser).Update("avatar", avatar)
	
	return nil
}