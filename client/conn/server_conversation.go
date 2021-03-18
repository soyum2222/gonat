package conn

import (
	"github.com/soyum2222/slog"
	"gonat/interface"
	"gonat/proto"
	"net"
	"sync"
)

type serverConversation struct {
	id            uint32
	remoteConn    net.Conn
	serverConn    net.Conn
	closeMu       sync.Mutex
	closeChan     chan struct{}
	cryptoHandler _interface.Safe
}

func (sc *serverConversation) Heartbeat() {
	//panic("implement me")
}

func (sc *serverConversation) Monitor() {

	data := make([]byte, 1024, 1024)

	for {

		select {

		case <-sc.closeChan:
			return

		default:

			n, err := sc.serverConn.Read(data)
			if err != nil {

				p := proto.Proto{Kind: proto.TCP_CLOSE_CONN, ConversationID: sc.id}
				_, err := sc.remoteConn.Write(p.Marshal(sc.cryptoHandler))
				if err != nil {
					slog.Logger.Error(err)
					sc.remoteConn.Close()
				}

				sc.Close()
				return
			}

			slog.Logger.Debug("server receive : ", string(data))
			slog.Logger.Debug("server receive len : ", n)

			p := proto.Proto{Kind: proto.TCP_COMM, ConversationID: sc.id, Body: data[0:n]}
			_, err = sc.remoteConn.Write(p.Marshal(sc.cryptoHandler))
			if err != nil {
				slog.Logger.Error(err)
				sc.Close()
				return
			}

		}

	}

}

func (sc *serverConversation) Close() {
	sc.serverConn.Close()

	sc.closeMu.Lock()
	defer sc.closeMu.Unlock()
	select {
	case _, ok := <-sc.closeChan:
		if !ok {
			return
		}
	default:
		close(sc.closeChan)

	}
}

func (sc *serverConversation) Send(b []byte) error {
	_, err := sc.serverConn.Write(b)
	return err
}
