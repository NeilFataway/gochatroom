package ws

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gochatroom/services"
)

var upgrader = websocket.Upgrader{}

func chat(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	defer func() {
		_ = conn.Close()
	}()

	if err != nil {
		log.WithError(err).Error("升级websocket连接失败")
		return
	}

	// 聊天室注册会话
	session, err := register(c, conn)
	if err != nil {
		log.WithError(err).Error("聊天室注册会话失败")
		return
	}

	isSessionClosed := false
	for {
		if isSessionClosed {
			break
		}
		message, err := session.Receive()
		if err != nil {
			var closeError *websocket.CloseError
			if errors.As(err, &closeError) || errors.As(err, &services.ClosedMsgReceived) {
				log.WithError(err).Infof("收到对端会话关闭报文: %s", session.GetRemoteIP())
			} else {
				log.WithError(err).Error("对端接受报文失败，退出会话")
			}
			_ = session.RemoveSession(session)
			isSessionClosed = true
		}
		log.WithField("message", message).Debugf("Receive message from peer: %s", session.GetRemoteIP())

		// 房间广播
		err = session.Room.Broadcast(message)
		if err != nil {
			log.WithError(err).Error("房间内广播失败，有终端未能发出信息")
		}
	}
}

func register(c *gin.Context, conn *websocket.Conn) (*services.Session, error) {
	//注册会话，进入聊天室

	// 校验房间ID
	roomId := c.Query("roomId")
	room, err := services.GetRoom(roomId)
	if err != nil {
		return nil, err
	}

	// 校验用户ID
	userId := c.Query("userId")
	user, err := services.GetUser(userId)
	if err != nil {
		return nil, err
	}

	// 创建Session
	session := &services.Session{
		Room: room,
		User: user,
		Peer: &WebsocketPeer{
			conn: conn,
		},
	}

	// session加入聊天室
	err = room.JoinSession(session)
	if err != nil {
		return nil, err
	}

	_ = room.Broadcast(services.FormLoginMsg(user))
	return session, nil
}
