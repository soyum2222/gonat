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
	remote_addr  *widget.Entry
	porxied_addr *widget.Entry
	dest_port    *widget.Entry
	crypt        *widget.Select
	crypt_key    *widget.Entry
	w            fyne.Window
	start        *widget.Button
)

func main() {

	app := app.New()

	config.GuiConfigPath = "./config.json"

	config.GuiLoad()

	w = app.NewWindow("GoNat")

	remote_addr = widget.NewEntry()
	porxied_addr = widget.NewEntry()
	dest_port = widget.NewEntry()
	crypt = widget.NewSelect([]string{"aes-128-cbc"}, nil)
	crypt_key = widget.NewPasswordEntry()

	remote_addr.SetText(config.CFG.RemoteAddr)
	porxied_addr.SetText(config.CFG.ProxiedAddr)
	dest_port.SetText(strconv.Itoa(config.CFG.DestPort))
	crypt.SetSelected(config.CFG.Crypt)
	crypt_key.SetText(config.CFG.CryptKey)

	start = widget.NewButton("Start", Strat)

	w.SetMainMenu(fyne.NewMainMenu(fyne.NewMenu("File",
		fyne.NewMenuItem("Start", Strat),
	)))

	form := &widget.Form{}

	form.Append("remote_addr", remote_addr)
	form.Append("porxied_addr", porxied_addr)
	form.Append("dest_port", dest_port)
	form.Append("crypt", crypt)
	form.Append("crypt_key", crypt_key)

	main_box := widget.NewVBox(
		form,
		widget.NewCheck("debug", func(b bool) {
			config.CFG.Debug = b

		}),
		start,
		widget.NewButton("Quit", func() {
			app.Quit()
		}))

	w.SetContent(main_box)

	w.ShowAndRun()
}

func Strat() {

	config.CFG.RemoteAddr = remote_addr.Text

	config.CFG.DestPort, _ = strconv.Atoi(dest_port.Text)

	config.CFG.CryptKey = crypt_key.Text

	config.CFG.ProxiedAddr = porxied_addr.Text

	config.CFG.Crypt = crypt.Selected

	port, err := strconv.Atoi(dest_port.Text)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}
	config.CFG.DestPort = port

	cfg := slog.TestSLogConfig()
	cfg.Debug = config.CFG.Debug

	err = slog.DefaultNew(cfg)

	start.SetText("stop")

	stop_sig, cancel := context.WithCancel(context.Background())
	start.OnTapped = func() {
		cancel()
	}

	conn.GuiStart(stop_sig, w)

	start.SetText("Start")
	start.OnTapped = Strat
}
