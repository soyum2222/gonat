package main

import (
	"github.com/soyum2222/slog"
	"gonat/client/config"
	"gonat/client/conn"
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"
)

func main() {

	config.Load()
	var err error

	cfg := slog.TestSLogConfig()
	cfg.Debug = config.CFG.Debug
	cfg.LogPath = config.CFG.LogPath
	cfg.LogFileName = "gonat"

	err = slog.DefaultNew(cfg)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := recover(); err != nil {
			slog.Logger.Panic(string(debug.Stack()), err)
		}
	}()

	if config.CFG.Debug {
		go http.ListenAndServe(config.CFG.PprofAddr, nil)
	}

	conn.Start()
}
