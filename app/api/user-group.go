package api

import (
	"gin-chat-svc/app/service"
	"gin-chat-svc/pkg/common/request"
	"gin-chat-svc/pkg/common/response"
	"gin-chat-svc/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserAndGroupList(ctx *gin.Context) {
	uuid := ctx.Query("uuid")
	interacts, err := service.NewUserAndGroupService.GetInteracts(uuid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, response.SuccessMsg(interacts))
}

// get user or group information by username or groupname
func GetUserOrGroupByName(ctx *gin.Context) {
	name := ctx.Query("name")
	getInfo, err := service.NewUserAndGroupService.GetUserOrGroupByName(name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessMsg(getInfo))
}

func PostInteract(ctx *gin.Context) {
	var userInteractReq *request.AllInteractReq

	ctx.ShouldBindJSON(&userInteractReq)
	postedInteract, err := service.NewUserAndGroupService.PostInteract(userInteractReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessMsg(postedInteract))
}

func DeleteInteract(ctx *gin.Context) {
	var userInteractReq request.AllInteractReq

	err := ctx.BindQuery(&userInteractReq)
	if err != nil {
		logger.Logger.Error("api", logger.Any("bind quert error", err))
	}
	logger.Logger.Info("api", logger.Any("userInteractRequest params: ", userInteractReq))

	deletedInteract, err := service.NewUserAndGroupService.DeleteInteract(&userInteractReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessMsg(deletedInteract))
}
