package ws

import (
	"gin-chat-svc/config"
	"gin-chat-svc/pkg/common/constant"
	"gin-chat-svc/pkg/logger"
	"gin-chat-svc/pkg/protocol"
	"gin-chat-svc/utility/kafka"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn		*websocket.Conn
	Name		string
	Send		chan []byte
}

func (c *Client) ReadPump() {
	defer func() {
		MyServer.Unregister <- c
		c.Conn.Close()
	}()

	for {
		c.Conn.PingHandler()
		// the server read messages using c.Conn.ReadMessage() (receive message)
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			logger.Logger.Error("websocket", logger.Any("client read message error", err.Error()))
			MyServer.Unregister <- c
			c.Conn.Close()
			break
		}
		
		msg := &protocol.Message {}
		proto.Unmarshal(message, msg)
		
		// pong
		if msg.Type == constant.HEART_BEAT {
			pong := &protocol.Message {
				Content: 	constant.PONG,
				Type: 		constant.HEART_BEAT,
			}

			pongByte, err2 := proto.Marshal(pong)
			if nil != err2 {
				logger.Logger.Error("websocket", logger.Any("client marshal message error", err2.Error()))
			}

			c.Conn.WriteMessage(websocket.BinaryMessage, pongByte)
		} else {
			conf, _ := config.GetConfig()
			if conf.ChannelType == constant.KAFKA {
				kafka.Send(message)
			} else {
				MyServer.Broadcast <- message
			}
		}
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()

	for message := range c.Send {
		// writes the message from server to client (send message)
		c.Conn.WriteMessage(websocket.BinaryMessage, message)
	}
}