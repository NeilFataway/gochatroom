package services

type Peer interface {
	Send(message *Message) error
	ShutDown() error
	GetRemoteIP() string
}

type Session struct {
	userName string
	Peer
}
