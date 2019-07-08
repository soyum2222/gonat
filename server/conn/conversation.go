package conn

import (
	"encoding/binary"
	"github.com/soyum2222/slog"
	"gonat/interface"
	"gonat/proto"
	"io"
	"net"
)

type local_conversation struct {
	user_conversation_map map[uint32]_interface.Conversation
	user_listener         net.Listener
	local_conn            net.Conn
	close_chan            chan struct{}
}

func (lc *local_conversation) Send([]byte) error {
	panic("implement me")
}

func (lc *local_conversation) Monitor() {
	l := make([]byte, 4, 4)
	p := proto.Proto{}

	for {

		select {
		case <-lc.close_chan:
			return
		default:
			_, err := io.ReadFull(lc.local_conn, l)
			//_, err := lc.local_conn.Read(l)
			if err != nil {
				slog.Logger.Error(err)
				lc.Close()
				return
			}

			data_len := binary.BigEndian.Uint32(l)

			data := make([]byte, data_len, data_len)

			_, err = io.ReadFull(lc.local_conn, data)
			//_, err = lc.local_conn.Read(data)
			if err != nil {
				slog.Logger.Error(err)
				lc.Close()
				return
			}

			p.Unmarshal(data)

			switch p.Kind {

			case proto.TCP_CLOSE_CONN:
				lc.user_conversation_map[p.ConversationID].Close()
				continue

			case proto.TCP_DIAL_ERROR:
				lc.user_conversation_map[p.ConversationID].Close()
				continue

			case proto.TCP_COMM:
				err := lc.user_conversation_map[p.ConversationID].Send(p.Body)
				if err != nil {
					slog.Logger.Error(err)
					lc.Close()
					return
				}

				slog.Logger.Debug("send user :", string(p.Body))
				slog.Logger.Debug("send user len:", len(p.Body))

			default:
				slog.Logger.Error("bad message:", p.Body)

			}
		}
	}
}

func (lc *local_conversation) Close() {
	for _, v := range lc.user_conversation_map {
		v.Close()
	}
	lc.user_listener.Close()
	close(lc.close_chan)
	lc.local_conn.Close()
}

func start_conversation(user_listen net.Listener, local_con net.Conn) {

	lc := local_conversation{}
	lc.local_conn = local_con
	lc.close_chan = make(chan struct{}, 1)
	lc.user_conversation_map = make(map[uint32]_interface.Conversation)
	lc.user_listener = user_listen
	go lc.Monitor()

	for conversation_id := 0; ; conversation_id++ {

		select {

		case <-lc.close_chan:
			return

		default:
			user_con, err := user_listen.Accept()
			if err != nil {
				slog.Logger.Error(err)
				return
			}

			p := proto.Proto{proto.TCP_CREATE_CONN, uint32(conversation_id), nil}
			_, err = lc.local_conn.Write(p.Marshal())
			if err != nil {
				user_con.Close()
				slog.Logger.Error(err)
				return
			}
			uc := user_conversation{}
			uc.local_conn = local_con
			uc.close_chan = make(chan struct{}, 1)
			uc.id = uint32(conversation_id)
			uc.user_conn = user_con
			lc.user_conversation_map[uc.id] = &uc

			go uc.Monitor()

		}

	}
}
