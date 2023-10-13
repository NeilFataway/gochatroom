package services

import (
	"errors"
	"time"
)

type MessageType byte

const (
	TextMessageType    MessageType = '1'
	PictureMessageType             = '2'
	VideoMessageType               = '3'
	SystemMessageType              = '4'
)

type Message struct {
	messageType MessageType
	data        []byte
	timestamp   time.Time
}

var InValidMessageType = errors.New("消息类型解析失败")

type MessageEncoder int

var DefaultMessageEncoder MessageEncoder

func (encoder *MessageEncoder) Encode(message *Message) ([]byte, error) {
	data := make([]byte, len(message.data)+1)
	data[0] = byte(message.messageType)
	copy(message.data, data[1:])
	return data, nil
}

type MessageDecoder int

var DefaultMessageDecoder MessageDecoder

func (decoder *MessageDecoder) Decode(data []byte) (*Message, error) {
	message := &Message{
		timestamp: time.Now(),
	}
	messageType := MessageType(data[0])
	switch messageType {
	case TextMessageType:
	case PictureMessageType:
	case VideoMessageType:
	case SystemMessageType:
		{

		}
	default:
		return nil, InValidMessageType
	}
	message.messageType = messageType
	copy(data[1:], message.data)

	return message, nil
}
