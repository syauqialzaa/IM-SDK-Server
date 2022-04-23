package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Group struct {
	ID       	int32                 	`json:"id" gorm:"primarykey"`
	Uuid      	string                	`json:"uuid" gorm:"type:varchar(150);not null;unique_index:idx_uuid;comment:uuid"`
	CreatedAt 	time.Time             	`json:"createAt"`
	UpdatedAt 	time.Time             	`json:"updatedAt"`
	DeletedAt 	soft_delete.DeletedAt	`json:"deletedAt"`
	UserId    	int32                 	`json:"userId" gorm:"index;comment:owner ID"`
	Name      	string                	`json:"name" gorm:"type:varchar(150);comment:group name"`
	Notice    	string                	`json:"notice" gorm:"type:varchar(350);comment:group notice"`
}