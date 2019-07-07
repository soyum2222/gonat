package main

import (
	"github.com/soyum2222/slog"
	"gonat/client/config"
	"gonat/client/conn"
)

func main() {

	config.Load()
	var err error
	slog.Logger, err = slog.DefaultNew(slog.TestSLogConfig)
	if err != nil {
		panic(err)
	}

	conn.Start()
}
