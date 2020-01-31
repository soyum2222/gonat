package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
)

var CFG config

type config struct {
	Port     int    `json:"port"`
	Debug    bool   `json:"debug"`
	Crypt    string `json:"crypt"`
	CryptKey string `json:"crypt_key"`
	LogPath  string `json:"log_path"`
}

func Load() {

	client_port := flag.Int("client_port", 0, "")
	debug := flag.Bool("debug", false, "")
	c := flag.String("c", "", "config file")
	crypt := flag.String("crypt", "", "crypt type")
	crypt_key := flag.String("crypt_key", "", "crypt key")
	log_path := flag.String("log_path", "", "log file path")

	flag.Parse()

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

		cfg.Port = *client_port
		cfg.Debug = *debug
		cfg.CryptKey = *crypt_key
		cfg.Crypt = *crypt
		cfg.LogPath = *log_path

	}

	CFG = cfg

}
