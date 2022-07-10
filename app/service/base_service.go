package service

import (
	"gin-chat-svc/utility/db"
	"time"
)

type UserService struct {}
type MessageService struct {}
type GroupService struct {}
type UserAndGroupService struct {}
type HttpUserService struct {}

var (
	Db 						= db.GetDB()
	NewUserService 			= new(UserService)
	NewMessageService 		= new(MessageService)
	NewGroupService 		= new(GroupService)
	NewUserAndGroupService 	= new(UserAndGroupService)
	NewHttpUserService 		= new(HttpUserService)

	timePast	= time.Date(
		2011, 12, 24, 10, 20, 0, 0, time.Local,
	)
)