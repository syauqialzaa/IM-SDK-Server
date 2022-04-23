package api

import (
	"gin-chat-svc/app/service"
	"gin-chat-svc/config"
	"gin-chat-svc/pkg/common/response"
	"gin-chat-svc/pkg/logger"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// the frontend gets the file stream by the file name and displays the file
func GetFile(ctx *gin.Context) {
	fileName := ctx.Param("fileName")
	logger.Logger.Info(fileName)

	conf, _ := config.GetConfig()
	data, _ := ioutil.ReadFile(conf.StaticFile + fileName)

	ctx.Writer.Write(data)
}

// upload avatar and other files
func SaveFile(ctx *gin.Context) {
	namePreffix := uuid.New().String()

	userUuid	:= ctx.PostForm("uuid")

	file, _		:= ctx.FormFile("file")
	fileName	:= file.Filename
	index		:= strings.LastIndex(fileName, ".")
	suffix		:= fileName[index:]

	newFileName := namePreffix + suffix
	conf, _		:= config.GetConfig()

	logger.Logger.Info("api", logger.Any("file name", conf.StaticFile + newFileName))
	logger.Logger.Info("api", logger.Any("userUuid name", userUuid))

	ctx.SaveUploadedFile(file, conf.StaticFile + newFileName)
	err := service.NewUserService.ModifyUserAvatar(newFileName, userUuid)
	
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
	}

	ctx.JSON(http.StatusOK, response.SuccessMsg(newFileName))
}