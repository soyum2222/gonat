package conn

import (
	"encoding/binary"
	"github.com/soyum2222/slog"
	"gonat/client/config"
	"gonat/proto"
	"gonat/safe"
	"gonat/sign"
	"net"
	"time"
)

func Start() {

	for {
		remote_conn, err := net.Dial("tcp", config.CFG.RemoteIp)
		if err != nil {
			slog.Logger.Error(err)
			time.Sleep(5 * time.Second)
			continue
		}

		start_conversation(remote_conn)
	}
}

func start_conversation(remote_conn net.Conn) {

	rc := remote_conversation{}
	rc.close_chan = make(chan struct{}, 1)
	rc.remote_conn = remote_conn
	rc.server_conversation_map.Init()
	rc.crypto_handler = safe.GetSafe(config.CFG.Crypt, config.CFG.CryptKey)

	port := make([]byte, 4, 4)

	binary.BigEndian.PutUint32(port, uint32(config.CFG.RemotePort))

	p := proto.Proto{
		Kind:           proto.TCP_SEND_PROTO,
		ConversationID: 0,
		Body:           sign.Signature(port),
	}
	_, err := rc.remote_conn.Write(p.Marshal(rc.crypto_handler))
	if err != nil {
		slog.Logger.Error(err)
		return
	}

	go rc.Heartbeat()
	rc.Monitor()

}
