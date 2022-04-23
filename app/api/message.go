package api

import (
	"gin-chat-svc/app/service"
	"gin-chat-svc/pkg/common/request"
	"gin-chat-svc/pkg/common/response"
	"gin-chat-svc/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// get message list
func GetMessage(ctx *gin.Context) {
	logger.Logger.Info(ctx.Query("uuid"))
	var messageRequest request.MessageRequest
	
	err := ctx.BindQuery(&messageRequest)
	if err != nil {
		logger.Logger.Error("api", logger.Any("bindQueryError", err))
	}
	logger.Logger.Info("api", logger.Any("messageRequest params: ", messageRequest))

	messages, err := service.NewMessageService.GetMessages(messageRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
	}

	ctx.JSON(http.StatusOK, response.SuccessMsg(messages))
}