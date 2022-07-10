package model

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID          	int32                 	`json:"id" gorm:"primarykey"`
	CreatedAt   	time.Time             	`json:"createAt" gorm:"DEFAULT:current_timestamp"`
	UpdatedAt   	time.Time             	`json:"updatedAt" gorm:"DEFAULT:null"`
	DeletedAt   	gorm.DeletedAt			`json:"deletedAt" gorm:"DEFAULT:null"`
	FromUuid		string					`json:"fromUuid" gorm:"type:varchar(150);not null;unique_index:idx_uuid;comment:fromUuid"`
	TargetUuid		string					`json:"targetUuid" gorm:"type:varchar(150);not null;unique_index:idx_uuid;comment:targetUuid"`
	Content     	string                	`json:"content" gorm:"type:varchar(2500)"`
	MessageType 	int16                 	`json:"messageType" gorm:"comment:[1]single chat or [2]group chat"`
	ContentType 	int16                 	`json:"contentType" gorm:"comment:text, common file, picture, audio, video, voice chat, video chat"`
	Pic         	string                	`json:"pic" gorm:"type:text;comment:thumbnail"`
	Url         	string                	`json:"url" gorm:"type:varchar(350);comment:file or image url"`
	SendTime		time.Time				`json:"sendTime"` // properly using time.Time
	SentTime		time.Time				`json:"sentTime"`
	ReadTime		time.Time				`json:"readTime"`
	DeleteFor		string					`json:"deleteFor"`
}

type MessageVisibility struct {
	ID				int32					`json:"id" gorm:"primarykey"`
	CreatedAt		time.Time				`json:"createdAt" gorm:"DEFAULT:current_timestamp"`
	MessageId		int32					`json:"messageId"`
	FromUuid		string					`json:"fromUuid" gorm:"type:varchar(150);not null;unique_index:idx_uuid;comment:fromUuid"`
	TargetUuid		string					`json:"targetUuid" gorm:"type:varchar(150);not null;unique_index:idx_uuid;comment:targetUuid"`
	WhosDelete		string					`json:"whosDelete"`
	DeleteStatus	string					`json:"deleteStatus" gorm:"comment:self/all"`
}

// func (m *Message) BeforeUpdate(tx *gorm.DB) error {
// 	tx.Statement.SetColumn("updated_at", timeNow)
// 	return nil
// }