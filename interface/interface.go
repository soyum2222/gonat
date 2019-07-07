package _interface

type Conversation interface {
	Monitor()
	Close()
	Send([]byte) error
}
