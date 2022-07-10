package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       	int32      			`json:"id" gorm:"primary_key;AUTO_INCREMENT;comment:id"`
	Uuid     	string     			`json:"uuid" gorm:"type:varchar(150);not null;unique_index:idx_uuid;comment:uuid"`
	CreatedAt 	time.Time  			`json:"createdAt" gorm:"DEFAULT:current_timestamp"`
	UpdatedAt 	time.Time 			`json:"updatedAt" gorm:"DEFAULT:null"`
	DeletedAt 	gorm.DeletedAt   	`json:"deletedAt" gorm:"DEFAULT:null"`
	Username 	string     			`json:"username" form:"username" binding:"required" gorm:"not null; comment:username"`
	Password 	string     			`json:"password" form:"password" binding:"required" gorm:"type:varchar(150);not null; comment:password"`
	Nickname 	string     			`json:"nickname" gorm:"comment:nickname"`
	Avatar   	string     			`json:"avatar" gorm:"type:varchar(150);comment:avatar"`
	Email    	string     			`json:"email" gorm:"type:varchar(80);column:email;comment:email"`
	StartOnline	time.Time			`json:"startOnline" gorm:"DEFAULT:null"`
	LastOnline	time.Time			`json:"lastOnline" gorm:"DEFAULT:null"`
}

// func (u *User) BeforeUpdate(tx *gorm.DB) error {
// 	tx.Statement.SetColumn("updated_at", timeNow)
// 	return nil
// }