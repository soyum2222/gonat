package conn

import (
	"encoding/binary"
	"github.com/soyum2222/slog"
	"gonat/common"
	"gonat/interface"
	"gonat/proto"
	"gonat/safe"
	"gonat/server/config"
	"gonat/sign"
	"io"
	"net"
	"strconv"
	"time"
)

type local_conversation struct {
	//map_mutex             sync.RWMutex
	//user_conversation_map map[uint32]_interface.Conversation //a user conn closed I didnt delete it,because I dont know how to do well . And I dont want use sync.Map :)
	user           common.ConversationTable
	user_listener  net.Listener
	local_conn     net.Conn
	close_chan     chan struct{}
	crypto_handler _interface.Safe
	timeout        *time.Timer
}

func (lc *local_conversation) Heartbeat() {
	panic("implement me")
}

func (lc *local_conversation) Send([]byte) error {
	panic("implement me")
}
func (lc *local_conversation) timeout_monitor() {

	for {
		select {
		case <-lc.timeout.C:
			slog.Logger.Info("user heartbeat timeout close the conn ", lc.local_conn.LocalAddr())
			lc.Close()
			return
		}
	}
}

//communication to gonat client
func (lc *local_conversation) Monitor() {
	l := make([]byte, 4, 4)
	p := proto.Proto{}

	lc.timeout = time.NewTimer(30 * time.Second)
	for {

		select {
		case <-lc.close_chan:
			return
		default:
			_, err := io.ReadFull(lc.local_conn, l)

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

			if conv, ok := lc.user.Load(p.ConversationID); conv == nil || !ok {
				continue
			}

			switch p.Kind {

			case proto.TCP_CLOSE_CONN:
				user_conn, _ := lc.user.Load(p.ConversationID)
				user_conn.Close()
				continue

			case proto.TCP_DIAL_ERROR:
				user_conn, _ := lc.user.Load(p.ConversationID)
				user_conn.Close()
				continue

			case proto.TCP_COMM:
				user_conn, _ := lc.user.Load(p.ConversationID)
				err := user_conn.Send(p.Body)
				if err != nil {
					//here no need return
					slog.Logger.Error(err)
					continue
					// 					lc.Close()
					// 					return
				}

				slog.Logger.Debug("send user :", string(p.Body))
				slog.Logger.Debug("send user len:", len(p.Body))

			case proto.Heartbeat:
				// 2019-12-24 find a bug , local to server conn disconnect , but server conn is normal
				// so need a conn heart beat timeout
				lc.timeout.Reset(30 * time.Second)
				continue
			default:
				slog.Logger.Error("bad message:", p.Body)

			}
		}
	}
}

func (lc *local_conversation) Close() {
	//for _, v := range lc.user_conversation_map {
	//	v.Close()
	//}
	lc.user.Range(func(key uint32, value _interface.Conversation) {
		value.Close()
	})

	err := lc.user_listener.Close()
	if err != nil {
		slog.Logger.Error("close user_listener error ", err, " listener info :", lc.user_listener.Addr())
	}

	select {
	case _, ok := <-lc.close_chan:
		if !ok {
			break
		}
	default:
		close(lc.close_chan)
	}

	lc.local_conn.Close()
}

func start_conversation(local_con net.Conn) {

	lc := local_conversation{}
	lc.crypto_handler = safe.GetSafe(config.CFG.Crypt, config.CFG.CryptKey)

	length := make([]byte, 4, 4)
	_, err := io.ReadFull(local_con, length)
	if err != nil {
		slog.Logger.Error("local conn read error , conn info : ", local_con.LocalAddr(), err)
		return
	}

	//signature verification
	ml := len((&proto.Proto{
		Kind:           proto.TCP_SEND_PROTO,
		ConversationID: 0,
		Body:           sign.Signature([]byte{0xff, 0xff, 0xff, 0xff}),
	}).Marshal(lc.crypto_handler))

	if binary.BigEndian.Uint32(length) > uint32(ml) {
		slog.Logger.Info("message is too long : ", local_con.RemoteAddr())
		return
	}

	data := make([]byte, binary.BigEndian.Uint32(length), binary.BigEndian.Uint32(length))
	_, err = io.ReadFull(local_con, data)
	if err != nil {
		slog.Logger.Error(err)
		return
	}

	p := proto.Proto{}
	p.Unmarshal(data, lc.crypto_handler)

	if !sign.Verifi(p.Body) || len(p.Body) != 8 { // one uint32 (port number) + 8 sign bit
		slog.Logger.Info("a bad message : ", local_con.RemoteAddr())
		return
	}

	port := binary.BigEndian.Uint32(p.Body[4:])

	listen, err := net.Listen("tcp", ":"+strconv.Itoa(int(port)))
	if err != nil {
		p := proto.Proto{Kind: proto.TCP_PORT_BIND_ERROR}
		_, _ = local_con.Write(p.Marshal(lc.crypto_handler))
		slog.Logger.Error(err)
		_ = local_con.Close()
		return
	}

	addr := listen.Addr().String()

	p = proto.Proto{proto.TCP_SEND_PROTO, 0, []byte(addr)}
	_, err = local_con.Write(p.Marshal(lc.crypto_handler))
	if err != nil {
		_ = local_con.Close()
		return
	}

	lc.local_conn = local_con
	lc.close_chan = make(chan struct{}, 1)
	lc.user.Init()
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
			uc.local = &lc
			uc.close_chan = make(chan struct{}, 1)
			uc.id = uint32(conversation_id)
			uc.user_conn = user_con
			uc.crypto_handler = lc.crypto_handler
			lc.user.Store(uc.id, &uc)
			go uc.Monitor()

		}

	}
}
