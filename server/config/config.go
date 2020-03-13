package config

import (
	"encoding/json"
	"flag"
	"fmt"
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
	PPROF    string `json:"pprof"`
	UIP      string `json:"uip"`
}

func Load() {

	p := flag.Int("p", 0, "")
	debug := flag.Bool("debug", false, "")
	c := flag.String("c", "", "config file")
	crypt := flag.String("ct", "", "crypt type")
	crypt_key := flag.String("k", "", "crypt key")
	log_path := flag.String("lp", "", "log file path")
	help := flag.Bool("help", false, "")
	uip := flag.String("uip", "", "")

	flag.Parse()

	if *help {

		fmt.Println("-p		listen port")
		fmt.Println("-debug	open or close debug mode")
		fmt.Println("-c		config file path")
		fmt.Println("-ct	crypt type now support aes-128-cbc")
		fmt.Println("-k		crypt password")
		fmt.Println("-lp	log file path if this value is null, no log file is create")
		fmt.Println("-pprof	when open debug mode, set pporf address eg:  127.0.0.1:8080")

		os.Exit(0)
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

		cfg.Port = *p
		cfg.Debug = *debug
		cfg.CryptKey = *crypt_key
		cfg.Crypt = *crypt
		cfg.LogPath = *log_path
		cfg.UIP = *uip

	}

	CFG = cfg

}
