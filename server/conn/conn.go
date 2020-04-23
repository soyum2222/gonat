package conn

import (
	"github.com/soyum2222/slog"
	"net"
	"sync"
)

var ClientTabel sync.Map
var UserTable sync.Map

func Start(port string) {

	slog.Logger.Info("server start")
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	for {
		defer func() {
			if err := recover(); err != nil {
				slog.Logger.Error(err)
			}
		}()

		local_con, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go start_conversation(local_con)
	}

}

//func local_task(local_con net.Conn) {
//
//	start_conversation(local_con)
//
//}
