package proto

import (
	"encoding/binary"
	"errors"
	"net"
)

const (
	TCP_CREATE_CONN = iota
	TCP_CLOSE_CONN
	TCP_COMM
	TCP_SEND_PROTO
	TCP_DIAL_ERROR
)

type Proto struct {
	Kind           uint32
	ConversationID uint32
	Body           []byte
}

func (p *Proto) Marshal() []byte {
	kind_b := make([]byte, 4, 4)
	id_b := make([]byte, 4, 4)
	len_b := make([]byte, 4, 4)
	binary.BigEndian.PutUint32(kind_b, p.Kind)
	binary.BigEndian.PutUint32(id_b, p.ConversationID)
	body := append(id_b, p.Body...)
	body = append(kind_b, body...)
	binary.BigEndian.PutUint32(len_b, uint32(len(body)))
	body = append(len_b, body...)

	return body
}

func (p *Proto) Unmarshal(b []byte) {

	kind_b := b[0:4]
	id_b := b[4:8]
	p.Body = b[8:]

	p.Kind = binary.BigEndian.Uint32(kind_b)
	p.ConversationID = binary.BigEndian.Uint32(id_b)

}

func GetPort(conn net.Conn) (uint32, error) {

	l := make([]byte, 8, 8)
	if _, err := conn.Read(l); err != nil {
		return 0, err
	}
	l_uint := binary.BigEndian.Uint64(l)
	b := make([]byte, l_uint, l_uint)

	if _, err := conn.Read(b); err != nil {
		return 0, err
	}

	if b[0] != TCP_SEND_PROTO {
		return 0, errors.New("msg not port")
	}

	return binary.BigEndian.Uint32(b[1:]), nil
}

func SendCreateConn() []byte {

	b := make([]byte, 9, 9)
	binary.BigEndian.PutUint64(b, 1)
	b[8] = TCP_CREATE_CONN
	return b
}
