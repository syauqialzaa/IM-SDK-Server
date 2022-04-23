package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type GroupMember struct {
	ID        	int32	               	`json:"id" gorm:"primarykey"`
	CreatedAt 	time.Time  	           	`json:"createAt"`
	UpdatedAt 	time.Time   	        `json:"updatedAt"`
	DeletedAt 	soft_delete.DeletedAt	`json:"deletedAt"`
	UserId    	int32                 	`json:"userId" gorm:"index;comment:user ID"`
	GroupId   	int32                 	`json:"groupId" gorm:"index;comment:group ID"`
	Nickname  	string                	`json:"nickname" gorm:"type:varchar(350);comment:nickname"`
	Mute      	int16                 	`json:"mute" gorm:"comment:mute"`
}