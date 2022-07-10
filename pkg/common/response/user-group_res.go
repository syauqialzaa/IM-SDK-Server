package response

type RespUserAndGroupList struct {
	ID			int32		`json:"id" gorm:"primarykey"`
	Uuid		string		`json:"uuid"`
	// users
	Username	string		`json:"username"`
	Avatar		string		`json:"avatar"`
	// groups
	Name		string		`json:"Name"`
}