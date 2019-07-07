package config

import "flag"

var Client_port int
var Debug bool

func Load() {

	client_port := flag.Int("client_port", 0, "")
	debug := flag.Bool("debug", false, "")

	flag.Parse()

	Client_port = *client_port
	Debug = *debug
}
