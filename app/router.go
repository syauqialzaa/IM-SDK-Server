package app

import (
	"gin-chat-svc/app/access"
	"gin-chat-svc/app/api"
	"gin-chat-svc/app/ws"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	app := gin.Default()

	app.Use(access.Cors())
	app.Use(access.Recovery)

	// var websocket func(ctx *gin.Context)
	socket := ws.RunWebsocket

	router := app.Group("")
	{
		router.GET("/user", api.GetUserList)
		router.GET("/user/:uuid", api.GetUserDetails)
		router.GET("/user/name", api.GetUserOrGroupByName)
		router.POST("/user/register", api.Register)
		router.POST("/user/login", api.Login)
		router.PUT("/user/:uuid", api.ModifyUser)
		router.DELETE("/user/:uuid", api.DeleteUser)

		// ---------------------------------------------
		router.POST("/interact", api.AddUserInteract)	// only user-to-user
		router.DELETE("/interact", api.DeleteUserInteract)
		// ---------------------------------------------
		router.GET("/interact", api.GetUserAndGroupList)
		router.POST("/all-interact", api.PostInteract)
		router.DELETE("/all-interact", api.DeleteInteract)	// users and groups interact

		router.GET("/message", api.GetMessage)
		router.GET("/message/id", api.GetMessageById)
		router.PUT("/message/id", api.ModifyMessage)
		router.DELETE("/message/id/delete-all", api.DeleteMsgForAll)
		router.DELETE("/message/id/delete-self", api.DeleteMsgForSelf)

		router.GET("/file/:fileName", api.GetFile)
		router.POST("/file", api.SaveFile)

		router.GET("/group", api.GetGroups)
		router.GET("/group/info/:uuid", api.GetGroupInfo)
		router.POST("/group/:uuid", api.SaveGroup)
		router.POST("/group/join/:userUuid/:groupUuid", api.JoinGroup)
		router.GET("/group/user/:uuid", api.GetGroupUsers)
		router.DELETE("/group/leave/:userUuid/:groupUuid", api.LeaveGroup)

		router.GET("/socket.io", socket)
	}

	request := app.Group("/http-req")
	{
		request.GET("/get-user-data", api.GetUserData)
		request.POST("/post-user-data", api.StoreUserData)
	}

	return app
}