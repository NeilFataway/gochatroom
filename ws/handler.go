package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gochatroom/services"
)

var upgrader = websocket.Upgrader{}

func chat(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	} else {
		register()
	}
	defer func() {
		_ = conn.Close()
	}()

	session := &services.Session{}

	chatRoom := c.Param("chatRoom")
	if chatRoom == "" || services.IsChatRoomExists(chatRoom) {
	}

	if err != nil {
		panic(err)
	}
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.WithError(err).Error("Receive message from peer failed.")
			break
		}
		log.WithField("message", message).Debug("Receive message from peer.")
		switch mt {
		case websocket.BinaryMessage:
			broadcast(message)
		case websocket.CloseMessage:
			unRegister()
			break
		}
	}
}

func broadcast(msg []byte) {
	//群发消息
}

func unRegister() {
	//注销会话，退出聊天室
}

func register() {
	//注册会话，进入聊天室
}
