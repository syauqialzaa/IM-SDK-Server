package response

// the response user data following json data from REST API target
type ResponseUserData struct {
	Username 	string     	`json:"username"`
	Password 	string     	`json:"password"`
	Nickname 	string     	`json:"nickname"`
	Avatar   	string     	`json:"avatar"`
	Email    	string     	`json:"email"`
}