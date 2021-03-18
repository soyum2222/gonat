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

type localConversation struct {
	//map_mutex             sync.RWMutex
	//user_conversation_map map[uint32]_interface.Conversation //a user conn closed I didnt delete it,because I dont know how to do well . And I dont want use sync.Map :)
	user          common.ConversationTable
	userListener  net.Listener
	localConn     net.Conn
	closeChan     chan struct{}
	cryptoHandler _interface.Safe
	timeout       *time.Timer
}

func (lc *localConversation) Heartbeat() {
	panic("implement me")
}

func (lc *localConversation) Send([]byte) error {
	panic("implement me")
}
func (lc *localConversation) timeoutMonitor() {

	for {
		select {
		case <-lc.timeout.C:
			slog.Logger.Info("user heartbeat timeout close the conn ", lc.localConn.LocalAddr())
			lc.Close()
			return
		}
	}
}

//communication to gonat client
func (lc *localConversation) Monitor() {
	l := make([]byte, 4, 4)
	p := proto.Proto{}

	lc.timeout = time.NewTimer(30 * time.Second)
	for {

		select {
		case <-lc.closeChan:
			return
		default:
			_, err := io.ReadFull(lc.localConn, l)

			if err != nil {
				slog.Logger.Error("local conn read error , conn info :", lc.localConn.LocalAddr(), err)
				lc.Close()
				return
			}

			dataLen := binary.BigEndian.Uint32(l)

			data := make([]byte, dataLen, dataLen)

			_, err = io.ReadFull(lc.localConn, data)
			//_, err = lc.localConn.Read(data)
			if err != nil {
				slog.Logger.Error(err)
				lc.Close()
				return
			}

			p.Unmarshal(data, lc.cryptoHandler)

			if conv, ok := lc.user.Load(p.ConversationID); conv == nil || !ok {
				continue
			}

			switch p.Kind {

			case proto.TCP_CLOSE_CONN:
				userConn, _ := lc.user.Load(p.ConversationID)
				userConn.Close()
				continue

			case proto.TCP_DIAL_ERROR:
				userConn, _ := lc.user.Load(p.ConversationID)
				userConn.Close()
				continue

			case proto.TCP_COMM:
				userConn, _ := lc.user.Load(p.ConversationID)
				err := userConn.Send(p.Body)
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

func (lc *localConversation) Close() {
	//for _, v := range lc.user_conversation_map {
	//	v.Close()
	//}
	lc.user.Range(func(key uint32, value _interface.Conversation) {
		value.Close()
	})

	ClientTable.Delete(lc.localConn.RemoteAddr().String())

	err := lc.userListener.Close()
	if err != nil {
		slog.Logger.Error("close userListener error ", err, " listener info :", lc.userListener.Addr())
	}

	select {
	case _, ok := <-lc.closeChan:
		if !ok {
			break
		}
	default:
		close(lc.closeChan)
	}

	lc.localConn.Close()
}

func startConversation(localCon net.Conn) {

	lc := localConversation{}
	lc.cryptoHandler = safe.GetSafe(config.CFG.Crypt, config.CFG.CryptKey)

	length := make([]byte, 4, 4)
	_, err := io.ReadFull(localCon, length)
	if err != nil {
		slog.Logger.Error("local conn read error , conn info : ", localCon.LocalAddr(), err)
		return
	}

	//signature verification
	ml := len((&proto.Proto{
		Kind:           proto.TCP_SEND_PROTO,
		ConversationID: 0,
		Body:           sign.Signature([]byte{0xff, 0xff, 0xff, 0xff}),
	}).Marshal(lc.cryptoHandler))

	if binary.BigEndian.Uint32(length) != uint32(ml)-4 {
		slog.Logger.Info("a bad message : ", localCon.RemoteAddr())
		return
	}

	data := make([]byte, binary.BigEndian.Uint32(length), binary.BigEndian.Uint32(length))
	_, err = io.ReadFull(localCon, data)
	if err != nil {
		slog.Logger.Error(err)
		return
	}

	p := proto.Proto{}
	p.Unmarshal(data, lc.cryptoHandler)

	if !sign.Verify(p.Body) || len(p.Body) != 8 { // one uint32 (port number) + 8 sign bit
		slog.Logger.Info("a bad message : ", localCon.RemoteAddr())
		return
	}

	port := binary.BigEndian.Uint32(p.Body[4:])

	listen, err := net.Listen("tcp", ":"+strconv.Itoa(int(port)))
	if err != nil {
		p := proto.Proto{Kind: proto.TCP_PORT_BIND_ERROR}
		_, _ = localCon.Write(p.Marshal(lc.cryptoHandler))
		slog.Logger.Error(err)
		_ = localCon.Close()
		return
	}

	addr := listen.Addr().String()

	// record
	ClientTable.Store(localCon.RemoteAddr().String(), addr)

	p = proto.Proto{Kind: proto.TCP_SEND_PROTO, Body: []byte(addr)}
	_, err = localCon.Write(p.Marshal(lc.cryptoHandler))
	if err != nil {
		_ = localCon.Close()
		return
	}

	lc.localConn = localCon
	lc.closeChan = make(chan struct{}, 1)
	lc.user.Init()
	lc.userListener = listen

	go lc.Monitor()

	for conversationId := 0; ; conversationId++ {

		select {

		case <-lc.closeChan:
			return

		default:
			userCon, err := listen.Accept()
			if err != nil {
				slog.Logger.Error(err)
				return
			}

			p := proto.Proto{Kind: proto.TCP_CREATE_CONN, ConversationID: uint32(conversationId)}
			_, err = lc.localConn.Write(p.Marshal(lc.cryptoHandler))
			if err != nil {
				userCon.Close()
				slog.Logger.Error(err)
				return
			}

			uc := user_conversation{}
			uc.local = &lc
			uc.close_chan = make(chan struct{}, 1)
			uc.id = uint32(conversationId)
			uc.user_conn = userCon
			uc.crypto_handler = lc.cryptoHandler
			lc.user.Store(uc.id, &uc)

			// recode
			UserTable.Store(userCon.RemoteAddr().String(), userCon.LocalAddr().String())

			go uc.Monitor()

		}

	}
}
