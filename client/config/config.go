package config

import (
	"flag"
)

var Remote_ip, Server_ip string
var Remote_port int

func Load() {

	server_ip := flag.String("server_ip", "127.0.0.1", "")

	remote_ip := flag.String("remote_ip", "", "")

	remote_l := flag.Int("remote_port", 0, "remote listen port")

	flag.Parse()

	Server_ip = *server_ip
	Remote_ip = *remote_ip
	Remote_port = *remote_l

}
