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
	err = slog.DefaultNew(func() slog.SLogConfig {
		cfg := slog.TestSLogConfig()
		cfg.Debug = config.CFG.Debug
		cfg.LogPath = config.CFG.LogPath
		return cfg
	})
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := recover(); err != nil {
			slog.Logger.Panic(string(debug.Stack()), err)
		}
	}()

	go http.ListenAndServe(":8808", nil)

	conn.Start()
}
