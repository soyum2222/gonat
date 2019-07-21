package main

import (
	"github.com/soyum2222/slog"
	"gonat/client/config"
	"gonat/client/conn"
)

func main() {

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

	defer func() {
		if err := recover(); err != nil {
			slog.Logger.Panic(err)
		}
	}()

	conn.Start()
}
