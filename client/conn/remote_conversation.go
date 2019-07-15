package conn

import (
	"encoding/binary"
	"github.com/soyum2222/slog"
	"gonat/client/config"
	"gonat/interface"
	"gonat/proto"
	"io"
	"net"
)

type remote_conversation struct {
	remote_conn             net.Conn
	server_conversation_map map[uint32]_interface.Conversation
	close_chan              chan struct{}
}

func (rc *remote_conversation) Monitor() {
	l := make([]byte, 4, 4)
	p := proto.Proto{}

	for {

		select {

		case <-rc.close_chan:
			return

		default:

			_, err := io.ReadFull(rc.remote_conn, l)
			//_, err := .Read(l)
			if err != nil {
				slog.Logger.Error(err)
				rc.Close()
				return
			}

			data_len := binary.BigEndian.Uint32(l)

			data := make([]byte, data_len, data_len)

			_, err = io.ReadFull(rc.remote_conn, data)
			//_, err = rc.remote_conn.Read(data)
			if err != nil {
				slog.Logger.Error(err)
				rc.Close()
				return
			}
			p.Unmarshal(data)

			switch p.Kind {

			case proto.TCP_CREATE_CONN:
				server_con, err := net.Dial("tcp", config.Server_ip)
				if err != nil {
					p.Kind = proto.TCP_DIAL_ERROR
					data := p.Marshal()
					err = rc.Send(data)
					if err != nil {
						rc.remote_conn.Close()
						close(rc.close_chan)
					}
				}
				sc := server_conversation{}
				sc.server_conn = server_con
				sc.remote_conn = rc.remote_conn
				sc.close_chan = make(chan struct{}, 1)
				sc.id = p.ConversationID
				go sc.Monitor()
				rc.server_conversation_map[p.ConversationID] = &sc

			case proto.TCP_COMM:
				err := rc.server_conversation_map[p.ConversationID].Send(p.Body)
				if err != nil {
					slog.Logger.Error(err)
					rc.server_conversation_map[p.ConversationID].Close()
					continue
				}

				slog.Logger.Debug("send server :", string(p.Body))
				slog.Logger.Debug("send server len:", len(p.Body))

			case proto.TCP_SEND_PROTO:
				slog.Logger.Info("remote port :", string(p.Body))

			}
		}

	}

}

func (rc *remote_conversation) Close() {
	for _, v := range rc.server_conversation_map {
		v.Close()
	}
	rc.remote_conn.Close()
	close(rc.close_chan)
}

func (rc *remote_conversation) Send(b []byte) error {
	_, err := rc.remote_conn.Write(b)
	return err
}
