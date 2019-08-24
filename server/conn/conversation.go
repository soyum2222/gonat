package conn

import (
	"encoding/binary"
	"github.com/soyum2222/slog"
	"gonat/interface"
	"gonat/proto"
	"gonat/safe"
	"gonat/server/config"
	"io"
	"net"
	"strconv"
	"strings"
)

type local_conversation struct {
	user_conversation_map map[uint32]_interface.Conversation //an user conn closed I didnt delete it,because I dont know how to do well . And I dont want use sync.Map :)
	user_listener         net.Listener
	local_conn            net.Conn
	close_chan            chan struct{}
	crypto_handler        _interface.Safe
}

func (lc *local_conversation) Heartbeat() {
	panic("implement me")
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
				slog.Logger.Error("local conn read error , conn info :", lc.local_conn.LocalAddr(), err)
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

			p.Unmarshal(data, lc.crypto_handler)

			if lc.user_conversation_map[p.ConversationID] == nil {
				continue
			}

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

			case proto.Heartbeat:
				continue
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
	err := lc.user_listener.Close()
	if err != nil {
		slog.Logger.Error("close user_listener error ", err, " listener info :", lc.user_listener.Addr())
	}
	close(lc.close_chan)
	lc.local_conn.Close()
}

func start_conversation(local_con net.Conn) {

	lc := local_conversation{}
	lc.crypto_handler = safe.GetSafe(config.Crypt, config.CryptKey)

	len_b := make([]byte, 4, 4)
	_, err := io.ReadFull(local_con, len_b)
	if err != nil {
		slog.Logger.Error("local conn read error , conn info :", local_con.LocalAddr(), err)
		return
	}

	if binary.BigEndian.Uint32(len_b) > 26 {
		slog.Logger.Info("the ip client is not gonat client", local_con.RemoteAddr())
		return
	}

	data_b := make([]byte, binary.BigEndian.Uint32(len_b), binary.BigEndian.Uint32(len_b))
	_, err = io.ReadFull(local_con, data_b)
	if err != nil {
		slog.Logger.Error(err)
		return
	}

	p := proto.Proto{}
	p.Unmarshal(data_b, lc.crypto_handler)

	if !strings.HasPrefix(string(p.Body), "gonat_port:") {
		return
	}

	port_b := p.Body[len([]byte("gonat_port:")):]

	port := binary.BigEndian.Uint32(port_b)

	listen, err := net.Listen("tcp", ":"+strconv.Itoa(int(port)))
	if err != nil {
		p := proto.Proto{Kind: proto.TCP_PORT_BIND_ERROR}
		local_con.Write(p.Marshal(lc.crypto_handler))
		slog.Logger.Error(err)
		local_con.Close()
		return
	}

	addr := listen.Addr().String()

	p = proto.Proto{proto.TCP_SEND_PROTO, 0, []byte(addr)}
	_, err = local_con.Write(p.Marshal(lc.crypto_handler))
	if err != nil {
		local_con.Close()
		return
	}

	lc.local_conn = local_con
	lc.close_chan = make(chan struct{}, 1)
	lc.user_conversation_map = make(map[uint32]_interface.Conversation)
	lc.user_listener = listen

	go lc.Monitor()

	for conversation_id := 0; ; conversation_id++ {

		select {

		case <-lc.close_chan:
			return

		default:
			user_con, err := listen.Accept()
			if err != nil {
				slog.Logger.Error(err)
				return
			}

			p := proto.Proto{proto.TCP_CREATE_CONN, uint32(conversation_id), nil}
			_, err = lc.local_conn.Write(p.Marshal(lc.crypto_handler))
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
			uc.crypto_handler = lc.crypto_handler
			lc.user_conversation_map[uc.id] = &uc

			go uc.Monitor()

		}

	}
}
