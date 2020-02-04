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
	server_ip   *widget.Entry
	remote_ip   *widget.Entry
	remote_port *widget.Entry
	crypt       *widget.Select
	crypt_key   *widget.Entry
	w           fyne.Window
	start       *widget.Button
)

func main() {

	app := app.New()

	config.GuiConfigPath = "./config.json"

	config.GuiLoad()

	w = app.NewWindow("GoNat")

	server_ip = widget.NewEntry()
	remote_ip = widget.NewEntry()
	remote_port = widget.NewEntry()
	crypt = widget.NewSelect([]string{"aes-128-cbc"}, nil)
	crypt_key = widget.NewPasswordEntry()

	server_ip.SetText(config.CFG.ServerIp)
	remote_ip.SetText(config.CFG.RemoteIp)
	remote_port.SetText(strconv.Itoa(config.CFG.RemotePort))
	crypt.SetSelected(config.CFG.Crypt)
	crypt_key.SetText(config.CFG.CryptKey)

	start = widget.NewButton("Star", Strat)

	w.SetMainMenu(fyne.NewMainMenu(fyne.NewMenu("File",
		fyne.NewMenuItem("Star", Strat),
	)))

	form := &widget.Form{}

	form.Append("server_ip", server_ip)
	form.Append("remote_ip", remote_ip)
	form.Append("remote_port", remote_port)
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

	config.CFG.RemotePort, _ = strconv.Atoi(remote_ip.Text)
	config.CFG.CryptKey = crypt_key.Text
	config.CFG.ServerIp = server_ip.Text
	config.CFG.Crypt = crypt.Selected
	port, err := strconv.Atoi(remote_port.Text)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}
	config.CFG.RemotePort = port

	err = slog.DefaultNew(func() slog.SLogConfig {
		cfg := slog.TestSLogConfig()
		cfg.Debug = config.CFG.Debug
		return cfg
	})

	start.SetText("stop")

	stop_sig, cancel := context.WithCancel(context.Background())
	start.OnTapped = func() {
		cancel()
	}

	conn.GuiStart(stop_sig, w)

	start.SetText("Start")
	start.OnTapped = Strat
}
