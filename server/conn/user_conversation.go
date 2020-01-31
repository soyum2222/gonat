package conn

import (
	"github.com/soyum2222/slog"
	"gonat/interface"
	"gonat/proto"
	"net"
	"sync"
)

type user_conversation struct {
	id             uint32
	user_conn      net.Conn
	local          *local_conversation
	close_mu       sync.Mutex
	close_chan     chan struct{}
	crypto_handler _interface.Safe
}

func (u *user_conversation) Heartbeat() {
	//panic("implement me")
}

func (u *user_conversation) Send(b []byte) error {
	_, err := u.user_conn.Write(b)
	return err
}

func (u *user_conversation) Monitor() {

	buf := make([]byte, 1024, 1024)
	for {

		select {
		case <-u.close_chan:
			return

		default:

			n, err := u.user_conn.Read(buf)
			if err != nil {
				slog.Logger.Info("a user close :", u.user_conn.RemoteAddr())
				p := proto.Proto{proto.TCP_CLOSE_CONN, u.id, nil}
				data := p.Marshal(u.crypto_handler)
				slog.Logger.Error(err)
				_, err := u.local.local_conn.Write(data)
				if err != nil {
					slog.Logger.Error(err)
					//u.local.Close()
				}
				u.Close()

				return
			}

			slog.Logger.Debug("user addr : ", u.user_conn.RemoteAddr(), " user receive : ", string(buf))
			slog.Logger.Debug("user receive len : ", n)

			err = u.send_to_local(buf[0:n])
			if err != nil {
				u.Close()
				slog.Logger.Error(err)
				return
			}

			slog.Logger.Debug("local send : ", string(buf))
			slog.Logger.Debug("local send len : ", n)
		}
	}

}

func (u *user_conversation) send_to_local(b []byte) error {

	p := proto.Proto{proto.TCP_COMM, u.id, b}
	_, err := u.local.local_conn.Write(p.Marshal(u.crypto_handler))
	return err

}

func (u *user_conversation) Close() {

	u.close_mu.Lock()
	defer u.close_mu.Unlock()
	u.user_conn.Close()
	u.local.user.Delete(u.id)

	select {
	case _, ok := <-u.close_chan:
		if !ok {
			return
		}
	default:
		close(u.close_chan)
	}

}
