package ws

import (
	"errors"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gochatroom/services"
	"strings"
	"time"
)

type WebsocketPeer struct {
	conn *websocket.Conn
}

func (wp *WebsocketPeer) Send(message *services.Message) error {
	encoder := services.DefaultMessageEncoder
	data, err := encoder.Encode(message)
	if err != nil {
		return err
	}

	if err = wp.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		return err
	}

	return nil
}

func (wp *WebsocketPeer) ShutDown() error {
	log.Infof("shutdown websocket connection from %s right now.", wp.conn.RemoteAddr().String())
	return wp.conn.Close()
}

func (wp *WebsocketPeer) GetRemoteIP() string {
	remoteAddr := wp.conn.RemoteAddr()

	if len(remoteAddr.String()) != 0 {
		return strings.Split(remoteAddr.String(), ":")[0]
	} else {
		return ""
	}
}

func (wp *WebsocketPeer) AbnormalShutDown(err error) error {
	message := websocket.FormatCloseMessage(websocket.CloseAbnormalClosure, err.Error())
	return wp.conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(100*time.Hour))
}

func (wp *WebsocketPeer) Receive() (*services.Message, error) {
	decoder := services.DefaultMessageDecoder
	mt, originalMsg, err := wp.conn.ReadMessage()
	if err != nil {
		var closeError *websocket.CloseError
		if errors.As(err, &closeError) {
			return nil, services.ClosedMsgReceived
		}
		return nil, err
	}
	switch mt {
	case websocket.BinaryMessage:
		return decoder.Decode(originalMsg)
	case websocket.TextMessage:
		return decoder.Decode(originalMsg)
	case websocket.CloseMessage:
		return nil, services.ClosedMsgReceived
	default:
		return nil, nil
	}
}
