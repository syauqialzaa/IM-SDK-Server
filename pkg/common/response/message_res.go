package response

import "time"

type MessageResponse struct {
	ID           	int32     	`json:"id" gorm:"primarykey"`
	FromUuid		string		`json:"fromUuid"`
	TargetUuid		string		`json:"targetUuid"`
	// FromUserId   	int32     	`json:"fromUserId" gorm:"index"`
	// ToUserId     	int32     	`json:"toUserId" gorm:"index"`
	Content      	string    	`json:"content" gorm:"type:varchar(2500)"`
	ContentType  	int16     	`json:"contentType" gorm:"comment:'text, voice, video'"`
	CreatedAt    	time.Time 	`json:"createAt"`
	FromUsername 	string    	`json:"fromUsername"`
	ToUsername   	string    	`json:"toUsername"`
	ToGroupName		string		`json:"toGroupName"`
	Avatar       	string    	`json:"avatar"`
	Url          	string    	`json:"url"`
}