package app

import (
	"gin-chat-svc/app/ws"
	"gin-chat-svc/config"
	"gin-chat-svc/pkg/common/constant"
	"gin-chat-svc/pkg/logger"
	"gin-chat-svc/utility/kafka"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RunApp() {
	gin.SetMode(gin.ReleaseMode)

	conf, _ := config.GetConfig()
	// db.GetDB()

	logger.InitLogger(conf.LogPath, conf.LogLevel)
	logger.Logger.Info("app", logger.Any("config", conf))

	if conf.ChannelType == constant.KAFKA {
		kafka.InitProducer(conf.KafkaTopic, conf.KafkaHost)
		kafka.InitConsumer(conf.KafkaHost)
		go kafka.ConsumerMsg(ws.ConsumerKafkaMsg)
	}

	logger.Logger.Info("app", logger.String("start", "start the server..."))

	handlers := Router()
	go ws.MyServer.StartServer()
	serve := &http.Server {
		Addr: 			":8000",
		Handler: 		handlers,
		ReadTimeout: 	10 * time.Second,
		WriteTimeout: 	10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := serve.ListenAndServe()
	if err != nil {
		logger.Logger.Error("app", logger.Any("server error", err))
	}
}