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

func DeleteUser(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	
	del := service.NewUserService.DeleteUser(uuid)
	ctx.JSON(http.StatusOK, response.SuccessMsg(del))
}

// func ModifyUserInfo(ctx *gin.Context) {
// 	var user model.User

// 	ctx.ShouldBindJSON(&user)
// 	logger.Logger.Debug("api", logger.Any("user", user))

// 	if err := service.NewUserService.ModifyUserInfo(&user); err != nil {
// 		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, response.SuccessMsg(nil))
// }

func ModifyUser(ctx *gin.Context) {
	var user *model.User
	var db = service.Db
	var uuid = ctx.Param("uuid")

	var NULL_ID int32 = 0
	db.First(&user, "uuid = ?", uuid)
	if NULL_ID == user.ID {
		ctx.JSON(http.StatusBadRequest, response.FailMsg("failed to get ID user."))
		return
	}

	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
		return
	}
	logger.Logger.Debug("api", logger.Any("user", user))

	service.NewUserService.ModifyUser(db, user)
	ctx.JSON(http.StatusOK, response.SuccessMsg(user))
}

func GetUserDetails(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	details := service.NewUserService.GetUserDetails(uuid)

	ctx.JSON(http.StatusOK, response.SuccessMsg(details))
}

func GetUserList(ctx *gin.Context) {
	uuid := ctx.Query("uuid")
	getList := service.NewUserService.GetUserList(uuid)

	ctx.JSON(http.StatusOK, response.SuccessMsg(getList))
}

func AddUserInteract(ctx *gin.Context) {
	var userInteractRequest request.InteractRequest

	ctx.ShouldBindJSON(&userInteractRequest)
	err := service.NewUserService.AddUserInteract(&userInteractRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessMsg(nil))
}

func DeleteUserInteract(ctx *gin.Context) {
	var userInteractRequest request.InteractRequest

	err := ctx.BindQuery(&userInteractRequest)
	if err != nil {
		logger.Logger.Error("api", logger.Any("bindQueryError", err))
	}
	logger.Logger.Info("api", logger.Any("userInteractRequest params: ", userInteractRequest))

	err = service.NewUserService.DeleteUserInteract(&userInteractRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessMsg(userInteractRequest))
}