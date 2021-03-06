package conn

import (
	"context"
	"encoding/binary"
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
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

		slog.Logger.Info("creating connection...")
		slog.Logger.Info("gonat remote address is : ", config.CFG.RemoteAddr)

		remoteConn, err := net.Dial("tcp", config.CFG.RemoteAddr)
		if err != nil {
			slog.Logger.Error(err)
			time.Sleep(5 * time.Second)
			continue
		}

		startConversation(remoteConn)
	}
}

func GuiStart(stopSignal context.Context, window fyne.Window) {

	content := window.Content()
	//box_v := *content.(*widget.Box).Children[0].(*widget.Form)

	temp := make([]fyne.CanvasObject, len(content.(*widget.Box).Children))
	copy(temp, content.(*widget.Box).Children)

	defer func() {
		content.(*widget.Box).Children[0] = temp[0]
	}()
	content.(*widget.Box).Children[0] = widget.NewLabelWithStyle("connecting...", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

label:
	//fmt.Println(config.Remote_ip)
	remoteConn, err := net.Dial("tcp", config.CFG.RemoteAddr)
	if err != nil {
		slog.Logger.Error(err)
		time.Sleep(5 * time.Second)

		select {
		case <-stopSignal.Done():
			remoteConn.Close()
		default:
			goto label
		}
	}
	content.(*widget.Box).Children[0] = widget.NewLabelWithStyle("connection succeeded", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-stopSignal.Done():
			remoteConn.Close()
			return
		}
	}()
	startConversation(remoteConn)
	cancel()

}

func startConversation(remoteConn net.Conn) {

	rc := remoteConversation{}
	rc.closeChan = make(chan struct{}, 1)
	rc.remoteConn = remoteConn
	rc.serverConversationMap.Init()
	rc.cryptoHandler = safe.GetSafe(config.CFG.Crypt, config.CFG.CryptKey)

	port := make([]byte, 4, 4)

	binary.BigEndian.PutUint32(port, uint32(config.CFG.DestPort))

	p := proto.Proto{
		Kind:           proto.TCP_SEND_PROTO,
		ConversationID: 0,
		Body:           sign.Signature(port),
	}
	_, err := rc.remoteConn.Write(p.Marshal(rc.cryptoHandler))
	if err != nil {
		slog.Logger.Error(err)
		return
	}

	slog.Logger.Info("connection established successfully")

	go rc.Heartbeat()
	rc.Monitor()

}
