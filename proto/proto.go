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
)

type Proto struct {
	Kind           uint32
	ConversationID uint32
	Body           []byte
}

func (p *Proto) Marshal(crypto_handler _interface.Safe) []byte {

	var err error
	p.Body, err = crypto_handler.Encrypt(p.Body)
	if err != nil {
		slog.Logger.Error(err)
		return nil
	}

	//len kind id body
	body := make([]byte, 12)
	binary.BigEndian.PutUint32(body[4:8], p.Kind)
	binary.BigEndian.PutUint32(body[8:12], p.ConversationID)
	binary.BigEndian.PutUint32(body[0:4], uint32(len(p.Body)+8))
	body = append(body, p.Body...)

	return body
}

func (p *Proto) Unmarshal(b []byte, crypto_handler _interface.Safe) {

	kind_b := b[0:4]
	id_b := b[4:8]
	p.Body = b[8:]
	var err error
	p.Body, err = crypto_handler.Decrypt(p.Body)
	if err != nil {
		slog.Logger.Error(err)
	}

	p.Kind = binary.BigEndian.Uint32(kind_b)
	p.ConversationID = binary.BigEndian.Uint32(id_b)

}
