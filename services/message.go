package services

import (
	"time"
)

type MessageType byte

const (
	TextMessageType             MessageType = '1'
	PictureMessageType                      = '2'
	VideoMessageType                        = '3'
	UserLoginSystemMessageType              = '4'
	UserLogoutSystemMessageType             = '5'
)

type Message struct {
	messageType MessageType
	fromWho     string
	data        []byte
	timestamp   time.Time
}

type MessageEncoder int

var DefaultMessageEncoder MessageEncoder

func (encoder *MessageEncoder) Encode(message *Message) ([]byte, error) {
	data := make([]byte, len(message.data)+1+36) //1为消息类型长度， 36为uuid长度
	data[0] = byte(message.messageType)
	copy(data[1:37], []byte(message.fromWho))
	copy(data[37:], message.data)
	return data, nil
}

type MessageDecoder int

var DefaultMessageDecoder MessageDecoder

func (decoder *MessageDecoder) Decode(data []byte) (*Message, error) {
	message := &Message{
		timestamp: time.Now(),
		data:      make([]byte, len(data)-1-36), //1为消息类型长度， 36为uuid长度
	}
	messageType := MessageType(data[0])
	switch messageType {
	case TextMessageType:
	case PictureMessageType:
	case VideoMessageType:
	case UserLoginSystemMessageType:
	case UserLogoutSystemMessageType:
	default:
		return nil, InValidMessageType
	}
	message.messageType = messageType
	message.fromWho = string(data[1:37])
	copy(message.data, data[37:])
	return message, nil
}

func FormLogoutMsg(user *User) *Message {

	return &Message{
		messageType: UserLogoutSystemMessageType,
		fromWho:     user.UserId,
		data:        []byte(user.UserName),
		timestamp:   time.Now(),
	}
}

func FormLoginMsg(user *User) *Message {
	return &Message{
		messageType: UserLoginSystemMessageType,
		fromWho:     user.UserId,
		data:        []byte(user.UserName),
		timestamp:   time.Now(),
	}
}
