package model

import (
	"time"

	"gorm.io/gorm"
	// "gorm.io/plugin/soft_delete"
)

type UserFriend struct {
	ID        	int32                 	`json:"id" gorm:"primarykey"`
	CreatedAt 	time.Time             	`json:"createAt"`
	UpdatedAt 	time.Time             	`json:"updatedAt"`
	DeletedAt 	gorm.DeletedAt 			`json:"deletedAt"`
	UserId    	int32                 	`json:"userId" gorm:"index;comment:user ID"`
	FriendId  	int32                 	`json:"friendId" gorm:"index;comment:friend ID"`
}