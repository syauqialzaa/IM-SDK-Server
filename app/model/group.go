package model

import (
	"time"

	"gorm.io/gorm"
)

type Group struct {
	ID       		int32				`json:"id" gorm:"primarykey"`
	Uuid      		string              `json:"uuid" gorm:"type:varchar(150);not null;unique_index:idx_uuid;comment:uuid"`
	CreatedAt 		time.Time           `json:"createdAt" gorm:"DEFAULT:current_timestamp"`
	UpdatedAt 		time.Time         	`json:"updatedAt" gorm:"DEFAULT:null"`
	DeletedAt 		gorm.DeletedAt		`json:"deletedAt" gorm:"DEFAULT:null"`
	CreatorUuid 	string              `json:"creatorUuid" gorm:"index;comment:creator uuid"`
	Creator			string				`json:"creator" gorm:"index;comment:creator"`
	Name      		string              `json:"name" gorm:"type:varchar(150);comment:group name"`
	Notice    		string              `json:"notice" gorm:"type:varchar(350);comment:group notice"`
}

// func (g *Group) BeforeUpdate(tx *gorm.DB) error {
// 	tx.Statement.SetColumn("updated_at", timeNow)
// 	return nil
// }