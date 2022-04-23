package api

import (
	"gin-chat-svc/app/model"
	"gin-chat-svc/app/service"
	"gin-chat-svc/pkg/common/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// get group list
func GetGroup(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	groups, err := service.NewGroupService.GetGroups(uuid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessMsg(groups))
}

// save group list
func SaveGroup(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	var group model.Group
	ctx.ShouldBindJSON(&group)

	service.NewGroupService.SaveGroup(uuid, group)
	ctx.JSON(http.StatusOK, response.SuccessMsg(nil))
}

// join group
func JoinGroup(ctx *gin.Context) {
	userUuid := ctx.Param("userUuid")
	groupUuid := ctx.Param("groupUuid")

	err := service.NewGroupService.JoinGroup(groupUuid, userUuid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessMsg(nil))
}

// get group member information
func GetGroupUsers(ctx *gin.Context) {
	groupUuid := ctx.Param("uuid")
	users := service.NewGroupService.GetUserIdByGroupUuid(groupUuid)
	
	ctx.JSON(http.StatusOK, response.SuccessMsg(users))
}