package model

import (
	"time"
)

type GroupMember struct {
	ID        	int32	               	`json:"id" gorm:"primarykey"`
	JoinedAt 	time.Time  	           	`json:"joinedAt" gorm:"DEFAULT:current_timestamp"`
	LeftAt	 	time.Time				`json:"leftAt" gorm:"DEFAULT:null"`
	GroupUuid   string                 	`json:"groupUuid"`
	Name		string					`json:"name" gorm:"type:varchar(150);comment:group name"`
	UserUuid    string                 	`json:"userUuid"`
	Username  	string                	`json:"username" gorm:"type:varchar(350);comment:member username"`
	Mute      	int16                 	`json:"mute" gorm:"comment:mute status"`
}

// func (gm *GroupMember) BeforeUpdate(tx *gorm.DB) error {
// 	tx.Statement.SetColumn("updated_at", timeNow)
// 	return nil
// }