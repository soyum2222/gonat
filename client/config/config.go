package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
)

var GuiConfigPath string
var CFG config

type config struct {
	Debug      bool   `json:"debug"`
	RemotePort int    `json:"remote_port"`
	LogPath    string `json:"log_path"`
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
	log_path := flag.String("log_path", "", "log file path")

	flag.Parse()

	if GuiConfigPath != "" {
		*c = GuiConfigPath
	}

	cfg := config{}
	if *c != "" {
		file, err := os.Open(*c)
		if err != nil {
			panic(err)
		}

		b, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(b, &cfg)
		if err != nil {
			panic(err)
		}

	} else {
		cfg.Crypt = *crypt
		cfg.ServerIp = *server_ip
		cfg.RemoteIp = *remote_ip
		cfg.RemotePort = *remote_l
		cfg.Debug = *debug
		cfg.CryptKey = *crypt_key
		cfg.LogPath = *log_path
	}

	CFG = cfg
}
func GuiLoad() {

	flag.Parse()

	if GuiConfigPath != "" {
		file, err := os.Open(GuiConfigPath)
		if err != nil {
			return
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

		CFG = cfg
	}
}
