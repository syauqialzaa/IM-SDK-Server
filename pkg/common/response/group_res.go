package response

import (
	"time"
)

type GroupResponse struct {
	Uuid      	string    	`json:"uuid"`
	GroupId   	int32     	`json:"groupId"`
	CreatedAt 	time.Time 	`json:"createdAt"`
	Name      	string    	`json:"name"`
	Notice    	string    	`json:"notice"`
}

// set group information =========
type GroupInfo struct {
	ID			int32			`json:"id"`
	Name		string			`json:"name"`
	Creator		string			`json:"creator"`
	CreatedAt	time.Time		`json:"createdAt"`
	Members		[]MemberList	`json:"members"`
}

type MemberList struct {
	Uuid		string		`json:"uuid"`
	Avatar		string		`json:"avatar"`
	Username	string		`json:"username"`
	Nickname	string		`json:"nickname"`
	Email		string		`json:"email"`
}
// ================================