package conn

import (
	"github.com/soyum2222/slog"
	"gonat/proto"
	"net"
)

type server_conversation struct {
	id          uint32
	remote_conn net.Conn
	server_conn net.Conn
	close_chan  chan struct{}
}

func (sc *server_conversation) Monitor() {

	data := make([]byte, 1024, 1024)

	for {

		select {

		case <-sc.close_chan:
			return

		default:

			n, err := sc.server_conn.Read(data)
			if err != nil {

				p := proto.Proto{proto.TCP_CLOSE_CONN, sc.id, nil}
				_, err := sc.remote_conn.Write(p.Marshal())
				if err != nil {
					slog.Logger.Error(err)
					sc.remote_conn.Close()
				}

				sc.Close()
				return
			}

			slog.Logger.Debug("server receive : ", string(data))
			slog.Logger.Debug("server receive len : ", n)

			p := proto.Proto{proto.TCP_COMM, sc.id, data[0:n]}
			_, err = sc.remote_conn.Write(p.Marshal())
			if err != nil {
				slog.Logger.Error(err)
				sc.Close()
				return
			}

		}

	}

}

func (sc *server_conversation) Close() {
	sc.server_conn.Close()

	select {
	case _, ok := <-sc.close_chan:
		if !ok {
			return
		}
	default:
		close(sc.close_chan)

	}
}

func (sc *server_conversation) Send(b []byte) error {
	_, err := sc.server_conn.Write(b)
	return err
}
