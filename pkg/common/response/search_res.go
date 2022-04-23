package response

import "gin-chat-svc/app/model"

type SearchResponse struct {
	User		model.User		`json:"user"`
	Group		model.Group		`json:"group"`
}