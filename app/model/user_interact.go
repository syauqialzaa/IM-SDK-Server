package model

import (
	"time"

)

// temp
type UserInteract struct {
	ID        			int32                 	`json:"id" gorm:"primarykey"`
	CreatedAt 			time.Time             	`json:"createdAt" gorm:"DEFAULT:current_timestamp"`
	// UpdatedAt 			*time.Time             	`json:"updatedAt"`
	DeletedAt 			time.Time	 			`json:"deletedAt" gorm:"DEFAULT:null"`
	UserId    			int32                 	`json:"userId" gorm:"index;comment:user ID"`
	Username			string					`json:"username"`
	InteractWithId 		int32                 	`json:"interactWithId"`
	InteractWith		string					`json:"interactWith"`
}

type AllUserInteract struct {
	ID					int32					`json:"id" gorm:"primarykey"`
	CreatedAt 			time.Time             	`json:"createdAt" gorm:"DEFAULT:current_timestamp"`
	DeletedAt 			time.Time	 			`json:"deletedAt" gorm:"DEFAULT:null"`
	UserUuid			string					`json:"userUuid"`
	Username			string					`json:"username"`
	InteractWithUuid	string					`json:"interactWithUuid"`
	InteractWith		string					`json:"interactWith"`
}

type UserInteractVisibility struct {
	ID					int32					`json:"id" gorm:"primarykey"`
	LastLeft			time.Time				`json:"lastLeft" gorm:"DEFAULT:null"`
	UserUuid			string					`json:"userUuid"`
	InteractWithUuid	string					`json:"interactWithUuid"`
	WhosLeave			string					`json:"whosLeave"`
	DeleteStatus		string					`json:"deleteStatus"`
}

// func (ui *UserInteract) BeforeUpdate(tx *gorm.DB) error {
// 	tx.Statement.SetColumn("update_at", timeNow)
// 	return nil
// }