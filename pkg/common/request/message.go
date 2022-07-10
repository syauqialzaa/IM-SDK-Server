package request

type MessageRequest struct {
	MessageType		int32		`json:"messageType"`
	Uuid			string		`json:"uuid"`
	InteractWith	string		`json:"interactWith"`
}

type MsgRequestById struct {
	ID				int32		`json:"id"`
	Uuid			string		`json:"uuid"`
	MessageType		int32		`json:"messageType"`
	InteractWith	string		`json:"interactWith"`
}