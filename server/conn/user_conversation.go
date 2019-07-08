package conn

import (
	"github.com/soyum2222/slog"
	"gonat/proto"
	"net"
)

type user_conversation struct {
	id         uint32
	user_conn  net.Conn
	local_conn net.Conn
	close_chan chan struct{}
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
				p := proto.Proto{proto.TCP_CLOSE_CONN, u.id, nil}
				data := p.Marshal()
				slog.Logger.Error(err)
				_, err := u.local_conn.Write(data)
				if err != nil {
					slog.Logger.Error(err)
					u.local_conn.Close()
				}
				u.Close()

				return
			}

			slog.Logger.Debug("user receive : ", string(buf))
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
	_, err := u.local_conn.Write(p.Marshal())
	return err

}

func (u *user_conversation) Close() {

	u.user_conn.Close()

	select {
	case _, ok := <-u.close_chan:
		if !ok {
			return
		}
	default:
		close(u.close_chan)
	}

}
