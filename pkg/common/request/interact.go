package request

type InteractRequest struct {
	Uuid				string
	InteractWith		string
}

type AllInteractReq struct {
	Uuid				string
	InteractWithUuid	string
}