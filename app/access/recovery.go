package access

import (
	"gin-chat-svc/pkg/common/response"
	"gin-chat-svc/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Recovery(ctx *gin.Context) {
	defer func() {
		if recover := recover(); recover != nil {
			logger.Logger.Error("access", logger.Any("gin catch error: ", recover))
			ctx.JSON(http.StatusInternalServerError, response.FailMsg("internal server error"))
		}
	}()

	ctx.Next()
}