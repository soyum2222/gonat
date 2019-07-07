package main

import (
	"github.com/soyum2222/slog"
	"gonat/server/config"
	"gonat/server/conn"
	"net/http"
	_ "net/http/pprof"
	"strconv"
)

func main() {
	go http.ListenAndServe(":8808", nil)
	config.Load()

	var err error
	slog.Logger, err = slog.DefaultNew(func() slog.SLogConfig {
		cfg := slog.TestSLogConfig()
		cfg.Debug = config.Debug
		return cfg
	})

	if err != nil {
		panic(err)
	}

	conn.Start(strconv.Itoa(config.Client_port))
}
