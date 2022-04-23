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
		router.PUT("/user", api.ModifyUserInfo)

		router.POST("/friend", api.AddFriend)
		// =====================================
		router.DELETE("/friend/:id", api.DeleteFriend)

		router.GET("/message", api.GetMessage)

		router.GET("/file/:fileName", api.GetFile)
		router.POST("/file", api.SaveFile)

		router.GET("/group/:uuid", api.GetGroup)
		router.POST("/group/:uuid", api.SaveGroup)
		router.POST("/group/join/:userUuid/:groupUuid", api.JoinGroup)
		router.GET("/group/user/:uuid", api.GetGroupUsers)

		router.GET("/socket.io", socket)
	}

	return app
}