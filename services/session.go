package services

import (
	"errors"
)

type Peer interface {
	Send(*Message) error
	Receive() (*Message, error)
	ShutDown() error
	GetRemoteIP() string
}

type Session struct {
	*Room
	*User
	Peer
}

func (s *Session) Receive() (*Message, error) {
	message, err := s.Peer.Receive()
	if errors.As(err, &ClosedMsgReceived) {
		return FormLogoutMsg(s.User), err
	} else {
		return message, err
	}
}
