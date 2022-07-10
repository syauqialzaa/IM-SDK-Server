package api

import (
	"gin-chat-svc/app/model"
	"gin-chat-svc/app/service"
	"gin-chat-svc/pkg/common/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// get group list
func GetGroups(ctx *gin.Context) {
	uuid := ctx.Query("uuid")

	groups, err := service.NewGroupService.GetGroups(uuid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessMsg(groups))
}

// get group info
func GetGroupInfo(ctx *gin.Context) {
	groupUuid := ctx.Param("uuid")
	
	groupInfo, err := service.NewGroupService.GetGroupInfo(groupUuid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessMsg(groupInfo))
}

// save group list
func SaveGroup(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	var group *model.Group
	ctx.ShouldBindJSON(&group)

	createdGroup, err := service.NewGroupService.SaveGroup(uuid, group)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, response.SuccessMsg(createdGroup))
}

// join group
func JoinGroup(ctx *gin.Context) {
	userUuid := ctx.Param("userUuid")
	groupUuid := ctx.Param("groupUuid")

	groupMember, err := service.NewGroupService.JoinGroup(groupUuid, userUuid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessMsg(groupMember))
}

func LeaveGroup(ctx *gin.Context) {
	userUuid := ctx.Param("userUuid")
	groupUuid := ctx.Param("groupUuid")

	grpInfo, err := service.NewGroupService.LeaveGroup(groupUuid, userUuid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessMsg(grpInfo))
}

// get group member information
func GetGroupUsers(ctx *gin.Context) {
	groupUuid := ctx.Param("uuid")
	users := service.NewGroupService.GetUserIdByGroupUuid(groupUuid)
	
	ctx.JSON(http.StatusOK, response.SuccessMsg(users))
}