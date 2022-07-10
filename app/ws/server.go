package ws

import (
	"encoding/base64"
	"gin-chat-svc/app/service"
	"gin-chat-svc/config"
	"gin-chat-svc/pkg/common/constant"
	"gin-chat-svc/pkg/common/request"
	"gin-chat-svc/pkg/common/suffix"
	"gin-chat-svc/pkg/logger"
	"gin-chat-svc/pkg/protocol"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
)

var MyServer = NewServer()

type Server struct {
	Clients		map[string] *Client
	Mutex		*sync.Mutex
	Broadcast	chan []byte
	Register	chan *Client
	Unregister	chan *Client
	Pagination	chan protocol.Message
}

func NewServer() *Server {
	return &Server {
		Mutex: 		&sync.Mutex {},
		Clients: 	make(map[string] *Client),
		Broadcast: 	make(chan []byte),
		Register: 	make(chan *Client),
		Unregister: make(chan *Client),
		Pagination: make(chan protocol.Message, 50),
	}
}

// consume the messages in kafka, and then put them directly into 
// the go channel for unified consumption
func ConsumerKafkaMsg(data []byte) {
	MyServer.Broadcast <- data
}

func (s *Server) StartServer() {
	logger.Logger.Info("websocket", logger.Any("server", "start server..."))

	for {
		select {
			case conn := <- s.Register:
				logger.Logger.Info("websocket", logger.Any("login", "new user has arrived: " + conn.Name))
				s.Clients[conn.Name] = conn
				
				// start online
				service.NewUserService.StartOnline(conn.Name)

				// autoload list of users and groups
				// service.NewUserAndGroupService.GetInteracts(conn.Name)

				msg := &protocol.Message {
					From: 		"system",
					To: 		conn.Name,
					Content: 	"welcome!",
				}

				protoMsg, _ := proto.Marshal(msg)
				conn.Send <- protoMsg

			case conn := <- s.Unregister:
				logger.Logger.Info("websocket", logger.Any("logout", conn.Name))

				if _, ok := s.Clients[conn.Name]; ok {
					close(conn.Send)
					delete(s.Clients, conn.Name)
				}

				// last online
				service.NewUserService.LastOnline(conn.Name)

			case message := <- s.Broadcast:

				// broadcast all messages in pagination in one block
				messages := make([]protocol.Message, 0)
				infLoop:
					for {
						select {
						case message := <- s.Pagination:
							messages = append(messages, message)
						default:
							break infLoop
						}
					}
					// if len(messages) > 0 {
					// 	for _, client := range s.Clients {
					// 		client.Send <- message
					// 	}
					// }

				msg := &protocol.Message {}
				proto.Unmarshal(message, msg)

				if msg.To != "" {
					// general messages, such as text messages, video file messages, etc.
					if msg.ContentType >= constant.TEXT && msg.ContentType <= constant.VIDEO {
						// saving messages will only be saved on one end of the socket, 
						// preventing the problem of message duplication after 
						// distributed deployment
						_, exist := s.Clients[msg.From]

						if exist {
							SaveMessage(msg)
						}

						if msg.MessageType == constant.MESSAGE_TYPE_USER {
							// get user information (details) when chatting
							service.NewUserService.GetUserDetails(msg.To)

							// get related messages
							msgReq := request.MessageRequest {
								MessageType: 	msg.MessageType,
								Uuid: 			msg.From,
								InteractWith: 	msg.To,
							}
							service.NewMessageService.GetMessages(msgReq)

							client, ok := s.Clients[msg.To]
							if ok {
								msgByte, err := proto.Marshal(msg)
								if err == nil {
									client.Send <- msgByte
								}
							}
						} else if msg.MessageType == constant.MESSAGE_TYPE_GROUP {
							SendGroupMessage(msg, s)
						}
					} else {
						// voice calls, video calls, etc., only support single-person chat, not group chat
						// do not save the file, forward it directly
						client, ok := s.Clients[msg.To]
						if ok {
							client.Send <- message
						}
					}
				} else if len(messages) > 0 {
					// there is no corresponding receiver to broadcast
					for id, conn := range s.Clients {
						logger.Logger.Info("websocket", logger.Any("allUser", id))

						select {
							case conn.Send <- message:
								default:
									close(conn.Send)
									delete(s.Clients, conn.Name)
						}
					}
				}
		}
	}
}

// to send a message to a group, you need to query all members of the group to send it in sequence
func SendGroupMessage(msg *protocol.Message, s *Server) {
	// get group information when in group chat
	// service.NewGroupService.GetGroupInfo(msg.To)

	// send a message to a group, find all users in the group and send it
	users := service.NewGroupService.GetUserIdByGroupUuid(msg.To)
	for _, user := range users {
		if user.Uuid == msg.From {
			continue
		}

		client, ok := s.Clients[user.Uuid]
		if !ok {
			continue
		}

		fromUserDetails := service.NewUserService.GetUserDetails(msg.From)
		// since when sending a group chat, from is an individual, and to is a group chat uuid. 
		// so when returning the message, change the form to the group chat uuid and 
		// unify it with the single chat
		msgSend := protocol.Message {
			Avatar: 		fromUserDetails.Avatar,
			FromUsername: 	msg.FromUsername,
			From: 			msg.To,
			To: 			msg.From,
			Content: 		msg.Content,
			ContentType: 	msg.ContentType,
			Type: 			msg.Type,
			MessageType: 	msg.MessageType,
			Url:		 	msg.Url,
		}

		msgByte, err := proto.Marshal(&msgSend)
		if err == nil {
			client.Send <- msgByte
		}
	}
}

// get all related messages with autoload
// func GetMessages(message *protocol.Message) {
// 	msg := request.MessageRequest {
// 		MessageType:	message.MessageType,
// 		Uuid: 			message.From,
// 		InteractWith: 	message.To,
// 	}
	
// 	if message.FromUsername != "" {
// 		service.NewMessageService.GetMessages(msg)
// 	}
// }

// save the message, if it is a text message, save it directly, if it is a file, 
// voice and other messages, after saving the file, save the corresponding file path
func SaveMessage(message *protocol.Message) {
	// if you upload a base64 string file, parse the file and save it
	if message.ContentType == 2 {
		url := uuid.New().String() + ".png"
		index := strings.Index(message.Content, "base64")
		index += 7

		content := message.Content
		content = content[index:]

		dataBuffer, dataErr := base64.StdEncoding.DecodeString(content)
		
		if dataErr != nil {
			logger.Logger.Error("websocket", logger.String("transfer base64 to file error", dataErr.Error()))
			return
		}

		conf, _ := config.GetConfig()
		err := ioutil.WriteFile(conf.StaticFile + url, dataBuffer, 0666)

		if err != nil {
			logger.Logger.Error("websocket", logger.String("write file error", err.Error()))
			return
		}

		message.Url = url
		message.Content = ""
	} else if message.ContentType == 3 {
		// ordinary file binary upload
		fileSuffix := suffix.GetFileType(message.File)
		var NULL_STRING string = ""

		if NULL_STRING == fileSuffix {
			fileSuffix = strings.ToLower(message.FileSuffix)
		}

		contentType := suffix.GetContentTypeBySuffix(fileSuffix)
		url := uuid.New().String() + "." + fileSuffix

		conf, _ := config.GetConfig()
		err := ioutil.WriteFile(conf.StaticFile + url, message.File, 0666)

		if err != nil {
			logger.Logger.Error("websocket", logger.String("write file error", err.Error()))
			return
		}

		message.Url			= url
		message.File		= nil
		message.ContentType	= contentType
	}

	MyServer.Pagination <- *message
	if message.FromUsername != "" {
		service.NewMessageService.SaveMessage(*message)
	}
}	