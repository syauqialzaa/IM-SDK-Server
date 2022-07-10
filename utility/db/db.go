package db

import (
	"fmt"
	"gin-chat-svc/app/model"
	"gin-chat-svc/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDB() *gorm.DB {
	defer func () {
		if recover := recover(); recover != nil {
			fmt.Println("recovered in f", recover)
		}
	} ()

	conf, _ := config.GetConfig()
	dsn := fmt.Sprint(conf.ElephantSQL)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config {
		Logger:		logger.Default.LogMode(logger.Info),
		NowFunc:	func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		panic("error to connect db " + err.Error())
	}

	sqlDB, _ := db.DB()

	db.AutoMigrate(
		model.GroupMember {},
		model.Group {},
		model.Message {},
		model.MessageVisibility {},
		model.AllUserInteract {},	// tmp from user_interacts
		model.UserInteract {},
		model.UserInteractVisibility {},
		model.User {},
	)

	// set db parameters
	// set the maximum number of connections in the database connection
	sqlDB.SetMaxOpenConns(100)
	// The maximum number of idle connections allowed by the connection.
	// If the number of connections that need to be executed for no sql task is greater
	// than 20, the excess connections will be closed by the connection.
	sqlDB.SetMaxIdleConns(20)

	return db
}