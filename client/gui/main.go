package main

import (
	"context"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
	"github.com/soyum2222/slog"
	"gonat/client/config"
	"gonat/client/conn"
	"strconv"
)

var (
	remoteAddr  *widget.Entry
	porxiedAddr *widget.Entry
	destPort    *widget.Entry
	crypt       *widget.Select
	cryptKey    *widget.Entry
	w           fyne.Window
	start       *widget.Button
)

func main() {

	app := app.New()

	config.GuiConfigPath = "./config.json"

	config.GuiLoad()

	w = app.NewWindow("GoNat")

	remoteAddr = widget.NewEntry()
	porxiedAddr = widget.NewEntry()
	destPort = widget.NewEntry()
	crypt = widget.NewSelect([]string{"aes-128-cbc"}, nil)
	cryptKey = widget.NewPasswordEntry()

	remoteAddr.SetText(config.CFG.RemoteAddr)
	porxiedAddr.SetText(config.CFG.ProxiedAddr)
	destPort.SetText(strconv.Itoa(config.CFG.DestPort))
	crypt.SetSelected(config.CFG.Crypt)
	cryptKey.SetText(config.CFG.CryptKey)

	start = widget.NewButton("Start", Start)

	w.SetMainMenu(fyne.NewMainMenu(fyne.NewMenu("File",
		fyne.NewMenuItem("Start", Start),
	)))

	form := &widget.Form{}

	form.Append("remote_addr", remoteAddr)
	form.Append("porxied_addr", porxiedAddr)
	form.Append("dest_port", destPort)
	form.Append("crypt", crypt)
	form.Append("crypt_key", cryptKey)

	mainBox := widget.NewVBox(
		form,
		widget.NewCheck("debug", func(b bool) {
			config.CFG.Debug = b

		}),
		start,
		widget.NewButton("Quit", func() {
			app.Quit()
		}))

	w.SetContent(mainBox)

	w.ShowAndRun()
}

func Start() {

	config.CFG.RemoteAddr = remoteAddr.Text

	config.CFG.DestPort, _ = strconv.Atoi(destPort.Text)

	config.CFG.CryptKey = cryptKey.Text

	config.CFG.ProxiedAddr = porxiedAddr.Text

	config.CFG.Crypt = crypt.Selected

	port, err := strconv.Atoi(destPort.Text)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}
	config.CFG.DestPort = port

	cfg := slog.TestSLogConfig()
	cfg.Debug = config.CFG.Debug

	err = slog.DefaultNew(cfg)

	start.SetText("stop")

	stopSig, cancel := context.WithCancel(context.Background())
	start.OnTapped = func() {
		cancel()
	}

	conn.GuiStart(stopSig, w)

	start.SetText("Start")
	start.OnTapped = Start
}
