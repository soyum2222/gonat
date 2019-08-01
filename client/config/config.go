package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
)

var Remote_ip, Server_ip, Crypt, CryptKey string
var Remote_port int
var Debug bool

type config struct {
	Debug      bool   `json:"debug"`
	RemotePort int    `json:"remote_port"`
	RemoteIp   string `json:"remote_ip"`
	ServerIp   string `json:"server_ip"`
	Crypt      string `json:"crypt"`
	CryptKey   string `json:"crypt_key"`
}

func Load() {

	server_ip := flag.String("server_ip", "127.0.0.1", "")

	remote_ip := flag.String("remote_ip", "", "")

	remote_l := flag.Int("remote_port", 0, "remote listen port")
	debug := flag.Bool("debug", false, "debug")
	c := flag.String("c", "", "config file")
	crypt := flag.String("crypt", "", "crypt type")
	crypt_key := flag.String("crypt_key", "", "crypt key")

	flag.Parse()

	if *c != "" {
		file, err := os.Open(*c)
		if err != nil {
			panic(err)
		}

		b, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
		cfg := config{}
		err = json.Unmarshal(b, &cfg)
		if err != nil {
			panic(err)
		}

		Remote_ip = cfg.RemoteIp
		Remote_port = cfg.RemotePort
		Server_ip = cfg.ServerIp
		Debug = cfg.Debug
		Crypt = cfg.Crypt
		CryptKey = cfg.CryptKey

	} else {
		Crypt = *crypt
		Server_ip = *server_ip
		Remote_ip = *remote_ip
		Remote_port = *remote_l
		Debug = *debug
		CryptKey = *crypt_key

	}
}
