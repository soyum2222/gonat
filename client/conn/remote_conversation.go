package conn

import (
	"encoding/binary"
	"github.com/soyum2222/slog"
	"gonat/client/config"
	"gonat/common"
	"gonat/interface"
	"gonat/proto"
	"io"
	"net"
	"sync"
	"time"
)

type remoteConversation struct {
	cryptoHandler         _interface.Safe
	remoteConn            net.Conn
	serverConversationMap common.ConversationTable // when keep long time gonat server conn this map  will leak memory
	closeChan             chan struct{}
	closeMu               sync.Mutex
}

func (rc *remoteConversation) Heartbeat() {

	for {

		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		select {
		case <-rc.closeChan:

			return

		case <-ticker.C:

			p := proto.Proto{
				Kind:           proto.Heartbeat,
				ConversationID: 0,
				Body:           make([]byte, 1, 1),
			}
			_, err := rc.remoteConn.Write(p.Marshal(rc.cryptoHandler))
			if err != nil {
				rc.Close()
				slog.Logger.Error(err)
				return
			}

		}

	}
}

func (rc *remoteConversation) Monitor() {
	l := make([]byte, 4, 4)
	p := proto.Proto{}

	for {

		select {

		case <-rc.closeChan:
			return

		default:

			_, err := io.ReadFull(rc.remoteConn, l)
			if err != nil {
				slog.Logger.Error(err)
				rc.Close()
				time.Sleep(time.Second * 2)
				return
			}

			dataLen := binary.BigEndian.Uint32(l)

			data := make([]byte, dataLen, dataLen)

			_, err = io.ReadFull(rc.remoteConn, data)
			if err != nil {
				slog.Logger.Error(err)
				rc.Close()
				return
			}

			p.Unmarshal(data, rc.cryptoHandler)

			switch p.Kind {

			case proto.TCP_CREATE_CONN:
				serverCon, err := net.Dial("tcp", config.CFG.ProxiedAddr)
				if err != nil {
					slog.Logger.Error(err)
					p.Kind = proto.TCP_DIAL_ERROR
					data := p.Marshal(rc.cryptoHandler)
					rc.Send(data)
					rc.remoteConn.Close()
					return
				}
				sc := serverConversation{}
				sc.serverConn = serverCon
				sc.remoteConn = rc.remoteConn
				sc.closeChan = make(chan struct{}, 1)
				sc.id = p.ConversationID
				sc.cryptoHandler = rc.cryptoHandler
				go sc.Monitor()
				rc.serverConversationMap.Store(p.ConversationID, &sc)

			case proto.TCP_COMM:
				// to server conversation
				scc, _ := rc.serverConversationMap.Load(p.ConversationID)
				err = scc.Send(p.Body)
				//err := rc.server_conversation_map[p.ConversationID].Send(p.Body)
				if err != nil {
					slog.Logger.Error(err)
					scc.Close()
					//rc.server_conversation_map[p.ConversationID].Close()
					continue
				}

				//slog.Logger.Debug("send server len:", len(p.Body))
				//slog.Logger.Debug("send server :", string(p.Body))

			case proto.TCP_SEND_PROTO:
				slog.Logger.Info("destination port :", string(p.Body))

			case proto.TCP_PORT_BIND_ERROR:
				slog.Logger.Info("remote port already bound please replace remote_port value")

			case proto.TCP_CLOSE_CONN:
				scc, _ := rc.serverConversationMap.Load(p.ConversationID)
				scc.Close()
				//rc.server_conversation_map[p.ConversationID].Close()
			}
		}

	}

}

func (rc *remoteConversation) Close() {
	//for _, v := range rc.server_conversation_map {
	//	v.Close()
	//}
	rc.serverConversationMap.Range(func(key uint32, value _interface.Conversation) {
		value.Close()
	})

	rc.closeMu.Lock()
	defer rc.closeMu.Unlock()
	rc.remoteConn.Close()

	select {
	case _, ok := <-rc.closeChan:
		if !ok {
			return
		}

	default:
		close(rc.closeChan)
	}
}

func (rc *remoteConversation) Send(b []byte) error {
	_, err := rc.remoteConn.Write(b)
	return err
}
