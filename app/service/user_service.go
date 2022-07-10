package service

import (
	"gin-chat-svc/app/model"
	"gin-chat-svc/pkg/common/constant"
	"gin-chat-svc/pkg/logger"
	"gin-chat-svc/pkg/misprint"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (u *UserService) Register(user *model.User) error {
	// db := db.GetDB()
	var userCount int64

	Db.First(&user, "username = ?", user.Username).Count(&userCount)
	if userCount > 0 {
		return misprint.New("user already exist")
	}

	// Db.Model(user).Where("username", user.Username).Count(&userCount)
	// if userCount > 0 {
	// 	return misprint.New("user already exist")
	// }

	// create unique user UUID
	user.Uuid = (constant.USER_TAG_UUID + uuid.New().String())
	// =====================================================================
	// user.CreatedAt = time.Now().Local()
	// user.DeletedAt = timeZero

	Db.Create(&user)

	return nil
}

func (u *UserService) Login(user *model.User) bool {
	var queryUser *model.User
	logger.Logger.Debug("service", logger.Any("user in service", user))

	Db.First(&queryUser, "username = ?", user.Username)
	logger.Logger.Debug("service", logger.Any("queryUser", queryUser))

	user.Uuid = queryUser.Uuid
	compare := queryUser.Password == user.Password
	
	return compare
}

func (u *UserService) DeleteUser(userUuid string) model.User {
	var queryUser *model.User

	Db.First(&queryUser, "uuid = ?", userUuid)
	if queryUser.ID == constant.NULL_ID {
		return *queryUser
	}

	Db.Delete(&queryUser)
	return *queryUser
}

// func (u *UserService) ModifyUserInfo(user *model.User) error {
// 	// db := db.GetDB()
// 	var queryUser *model.User

// 	Db.First(&queryUser, "username = ?", user.Username) // user.Username
// 	logger.Logger.Debug("service", logger.Any("queryUser", queryUser))

// 	var ID_NIL int32 = 0
// 	if ID_NIL == queryUser.ID {
// 		return misprint.New("user doesn't exist")
// 	}

// 	queryUser.Nickname	= user.Nickname
// 	queryUser.Email		= user.Email
// 	queryUser.Password	= user.Password

// 	Db.Save(&queryUser)

// 	return nil
// }

func (u *UserService) ModifyUser(db *gorm.DB, user *model.User) error {
	// var user model.User
	if err := db.Model(&user).Updates(map[string]interface{} {
		"username":		user.Username,
		"password":		user.Password,
		"nickname":		user.Nickname,
		"email":		user.Email,
	}).Error; err != nil {
		return misprint.New("failed to update.")
	}

	db.Model(&user).Update("updated_at", time.Now().Local())
	// db.Save(user)
	return nil
}

// update user start online in db
func (u *UserService) StartOnline(userUuid string) (time.Time, error) {
	var queryUser *model.User
	var timestamp time.Time

	Db.First(&queryUser, "uuid = ?", userUuid)
	if queryUser.ID == constant.NULL_ID {
		return timestamp, misprint.New("user doesn't exist.")
	}

	Db.Model(&queryUser).Update("start_online", time.Now().Local())
	// return time.Time in model to NULL
	Db.Raw(`UPDATE users SET last_online = NULL WHERE uuid = ?`, queryUser.Uuid).Scan(&queryUser)

	return timestamp, nil
}

// update user last online in db
func (u *UserService) LastOnline(userUuid string) (time.Time, error) {
	var queryUser *model.User
	var timestamp time.Time

	Db.First(&queryUser, "uuid = ?", userUuid)
	if queryUser.ID == constant.NULL_ID {
		return timestamp, misprint.New("user doesn't exist.")
	}

	Db.Model(&queryUser).Update("last_online", time.Now().Local())

	return timestamp, nil
}

func (u *UserService) GetUserDetails(uuid string) model.User {
	// db := db.GetDB()
	var queryUser *model.User

	Db.Select("uuid", "username", "nickname", "avatar", "email").
		First(&queryUser, "uuid = ?", uuid)

	return *queryUser
}

func (u *UserService) GetUserList(uuid string) []model.User {
	// db := db.GetDB()
	var queryUser *model.User

	Db.First(&queryUser, "uuid = ?", uuid)
	if queryUser.ID == constant.NULL_ID {
		return nil
	}

	var queryUsers []model.User

	Db.Raw(`
		SELECT u.username, u.uuid, u.avatar 
		FROM user_interacts AS ui JOIN users AS u ON ui.interact_with_id = u.id 
		WHERE ui.deleted_at IS NULL AND ui.user_id = ?
	`, queryUser.ID).Scan(&queryUsers)

	return queryUsers
}

// modify avatar
func (u *UserService) ModifyUserAvatar(avatar string, userUuid string) error {
	// db := db.GetDB()
	var queryUser *model.User

	Db.First(&queryUser, "uuid = ?", userUuid)
	if queryUser.ID == constant.NULL_ID {
		return misprint.New("user doesn't exist")
	}

	Db.Model(&queryUser).Update("avatar", avatar)

	return nil
}