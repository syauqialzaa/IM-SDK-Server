package api

import (
	"gin-chat-svc/app/model"
	"gin-chat-svc/app/service"
	"gin-chat-svc/pkg/common/constant"
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
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessMsg(messages))
}

func GetMessageById(ctx *gin.Context) {
	logger.Logger.Info(ctx.Query("uuid"))
	var getMsgByIdReq request.MsgRequestById
	// var message model.Message

	err := ctx.BindQuery(&getMsgByIdReq)
	if err != nil {
		logger.Logger.Error("api", logger.Any("bindQueryError", err))
	}
	logger.Logger.Info("api", logger.Any("messageRequest params: ", getMsgByIdReq))

	// ctx.ShouldBindJSON(&message)
	msg, err := service.NewMessageService.GetMessageById(getMsgByIdReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessMsg(msg))
}

// update message
func ModifyMessage(ctx *gin.Context) {
	logger.Logger.Info(ctx.Query("uuid"))

	var db = service.Db
	var modifyMsgReq request.MsgRequestById
	var message *model.Message
	var from *model.User

	err := ctx.BindQuery(&modifyMsgReq)
	if err != nil {
		logger.Logger.Error("api", logger.Any("bindQueryError", err))
	}
	logger.Logger.Info("api", logger.Any("messageRequest params: ", modifyMsgReq))

	db.First(&from, "uuid = ?", modifyMsgReq.Uuid)
	if from.ID <= 0 {
		ctx.JSON(http.StatusBadRequest, response.FailMsg("failed to get ID user."))
		return
	}
	
	db.First(&message, "id = ?", modifyMsgReq.ID)
	if message.ID <= 0 {
		ctx.JSON(http.StatusBadRequest, response.FailMsg("failed to get ID message."))
		return
	}

	db.First(&message, "message_type = ?", modifyMsgReq.MessageType)
	if modifyMsgReq.MessageType == constant.MESSAGE_TYPE_USER {
		var interactWith *model.User
		
		db.First(&interactWith, "uuid = ?", modifyMsgReq.InteractWith)
		if interactWith.ID <= 0 {
			ctx.JSON(http.StatusBadRequest, response.FailMsg("failed to get ID user."))
			return
		}

	} else if modifyMsgReq.MessageType == constant.MESSAGE_TYPE_GROUP {
		var interactGroup *model.Group
		
		db.First(&interactGroup, "uuid = ?", modifyMsgReq.InteractWith)
		if interactGroup.ID <= 0 {
			ctx.JSON(http.StatusBadRequest, response.FailMsg("failed to get ID group."))
			return
		}
	}

	err = ctx.ShouldBindJSON(&message)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
		return
	}
	logger.Logger.Debug("api", logger.Any("update message", message))

	service.NewMessageService.ModifyMessage(message)
	ctx.JSON(http.StatusOK, response.SuccessMsg(message))
}

func DeleteMsgForAll(ctx *gin.Context) {
	logger.Logger.Info(ctx.Query("uuid"))
	var getMsgByIdReq request.MsgRequestById

	err := ctx.BindQuery(&getMsgByIdReq)
	if err != nil {
		logger.Logger.Error("api", logger.Any("bindQueryError", err))
	}
	logger.Logger.Info("api", logger.Any("messageRequest params: ", getMsgByIdReq))

	err = service.NewMessageService.DeleteMsgForAll(getMsgByIdReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessMsg(getMsgByIdReq))
}

func DeleteMsgForSelf(ctx *gin.Context) {
	logger.Logger.Info(ctx.Query("uuid"))
	var getMsgByIdReq request.MsgRequestById

	err := ctx.BindQuery(&getMsgByIdReq)
	if err != nil {
		logger.Logger.Error("api", logger.Any("bindQueryError", err))
	}
	logger.Logger.Info("api", logger.Any("messageRequest params: ", getMsgByIdReq))

	err = service.NewMessageService.DeleteMsgForSelf(getMsgByIdReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailMsg(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessMsg(getMsgByIdReq))
}