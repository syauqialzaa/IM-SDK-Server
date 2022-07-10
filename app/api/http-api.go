package api

import (
	"gin-chat-svc/app/service"
	"gin-chat-svc/pkg/common/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserData(ctx *gin.Context) {
	var userData []response.ResponseUserData

	data, err := service.NewHttpUserService.GetUserData(userData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessMsg(data))
}

func StoreUserData(ctx *gin.Context) {
	var userData []response.ResponseUserData

	data, err := service.NewHttpUserService.StoreUserData(userData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessMsg(data))
}