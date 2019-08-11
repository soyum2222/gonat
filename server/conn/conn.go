package conn

import (
	"github.com/soyum2222/slog"
	"net"
)

func Start(port string) {

	err := slog.Logger.Info("server start")
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	for {
		local_con, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go local_task(local_con)
	}

}

func local_task(local_con net.Conn) {

	start_conversation(local_con)

}
