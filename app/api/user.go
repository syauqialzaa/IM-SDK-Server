package api

import (
	"gin-chat-svc/app/model"
	"gin-chat-svc/app/service"
	"gin-chat-svc/pkg/common/request"
	"gin-chat-svc/pkg/common/response"
	"gin-chat-svc/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	var user model.User

	ctx.ShouldBindJSON(&user)
	err := service.NewUserService.Register(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
	}

	ctx.JSON(http.StatusOK, response.SuccessMsg(user))
}

func Login(ctx *gin.Context) {
	var user model.User

	ctx.ShouldBindJSON(&user)
	logger.Logger.Debug("api", logger.Any("user", user))

	if service.NewUserService.Login(&user) {
		ctx.JSON(http.StatusOK, response.SuccessMsg(user))
		return
	}

	ctx.JSON(http.StatusBadRequest, response.FailMsg("login failed"))
}

func ModifyUserInfo(ctx *gin.Context) {
	var user model.User

	ctx.ShouldBindJSON(&user)
	logger.Logger.Debug("api", logger.Any("user", user))

	if err := service.NewUserService.ModifyUserInfo(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessMsg(nil))
}

func GetUserDetails(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	details := service.NewUserService.GetUserDetails(uuid)

	ctx.JSON(http.StatusOK, response.SuccessMsg(details))
}

// get user information by username
func GetUserOrGroupByName(ctx *gin.Context) {
	name := ctx.Query("name")
	getInfo := service.NewUserService.GetUserOrGroupByName(name)

	ctx.JSON(http.StatusOK, response.SuccessMsg(getInfo))
}

func GetUserList(ctx *gin.Context) {
	uuid := ctx.Query("uuid")
	getList := service.NewUserService.GetUserList(uuid)

	ctx.JSON(http.StatusOK, response.SuccessMsg(getList))
}

func AddFriend(ctx *gin.Context) {
	var userFriendRequest request.FriendRequest

	ctx.ShouldBindJSON(&userFriendRequest)
	err := service.NewUserService.AddFriend(&userFriendRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessMsg(nil))
}

func DeleteFriend(ctx *gin.Context) {
	id := ctx.Param("id")
	service.NewUserService.DeleteFriend(id)
	var updatedData model.UserFriend

	ctx.JSON(http.StatusOK, response.SuccessMsg(updatedData))
}