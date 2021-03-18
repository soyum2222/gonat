package proto

import (
	"encoding/binary"
	"github.com/soyum2222/slog"
	"gonat/interface"
)

const (
	TCP_CREATE_CONN = iota
	TCP_CLOSE_CONN
	TCP_COMM
	TCP_SEND_PROTO
	TCP_DIAL_ERROR
	TCP_PORT_BIND_ERROR
	Heartbeat
	BAD_MESSAGE
)

type Proto struct {
	Kind           uint32
	ConversationID uint32
	Body           []byte
}

func (p *Proto) Marshal(cryptoHandler _interface.Safe) []byte {

	var err error
	p.Body, err = cryptoHandler.Encrypt(p.Body)
	if err != nil {
		slog.Logger.Error(err)
		return nil
	}

	//len kind id body
	body := make([]byte, 12)
	binary.BigEndian.PutUint32(body[4:8], p.Kind)
	binary.BigEndian.PutUint32(body[8:12], p.ConversationID)
	binary.BigEndian.PutUint32(body[0:4], uint32(len(p.Body)+8)) // kind is 4 byte , conversationid is 4 byte .
	body = append(body, p.Body...)

	return body
}

func (p *Proto) Unmarshal(b []byte, cryptoHandler _interface.Safe) {
	if len(b) < 8 {
		p.Kind = BAD_MESSAGE
		return
	}
	var err error
	p.Kind = binary.BigEndian.Uint32(b[0:4])
	p.ConversationID = binary.BigEndian.Uint32(b[4:8])
	p.Body, err = cryptoHandler.Decrypt(b[8:])
	if err != nil {
		slog.Logger.Error(err)
	}
}
