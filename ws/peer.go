package ws

import (
	"github.com/gorilla/websocket"
	"gochatroom/services"
	"strings"
	"time"
)

type WebsocketPeer struct {
	conn websocket.Conn
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
