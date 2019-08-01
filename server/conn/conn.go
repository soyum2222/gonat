package conn

import (
	"encoding/binary"
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

	port_b := make([]byte, 4, 4)
	_, err := local_con.Read(port_b)
	if err != nil {
		slog.Logger.Error(err)
		return
	}
	port := binary.BigEndian.Uint32(port_b)

	slog.Logger.Info("client link port :", port)

	start_conversation(port, local_con)

}
