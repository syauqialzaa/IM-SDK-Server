package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Message struct {
	ID          	int32                 	`json:"id" gorm:"primarykey"`
	CreatedAt   	time.Time             	`json:"createAt"`
	UpdatedAt   	time.Time             	`json:"updatedAt"`
	DeletedAt   	soft_delete.DeletedAt 	`json:"deletedAt"`
	FromUserId  	int32                 	`json:"fromUserId" gorm:"index"`
	ToUserId    	int32                 	`json:"toUserId" gorm:"index;comment:user ID or group ID"`
	Content     	string                	`json:"content" gorm:"type:varchar(2500)"`
	MessageType 	int16                 	`json:"messageType" gorm:"comment:single chat or group chat"`
	ContentType 	int16                 	`json:"contentType" gorm:"comment:text, common file, picture, audio, video, voice chat, video chat"`
	Pic         	string                	`json:"pic" gorm:"type:text;comment:thumbnail"`
	Url         	string                	`json:"url" gorm:"type:varchar(350);comment:file or image url"`
}