package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var GuiConfigPath string
var CFG config

type config struct {
	Debug       bool   `json:"debug"`
	DestPort    int    `json:"dest_port"`
	LogPath     string `json:"log_path"`
	RemoteAddr  string `json:"remote_addr"`
	ProxiedAddr string `json:"proxied_addr"`
	Crypt       string `json:"crypt"`
	CryptKey    string `json:"crypt_key"`
	PprofAddr   string `json:"pprof_addr"`
}

func PrintHelp() {

	fmt.Println("-p         be proxied server address eg:  127.0.0.1:80")
	fmt.Println("-r         gonat server address ")
	fmt.Println("-dp        destination port this port will to gonat server listen")
	fmt.Println("-debug     open or close debug mode")
	fmt.Println("-c         config file path")
	fmt.Println("-ct        crypt type now support aes-128-cbc")
	fmt.Println("-k         crypt password")
	fmt.Println("-lp        log file path if this value is null, no log file is create")
	fmt.Println("-pprof     when open debug mode, set pporf address eg:  127.0.0.1:8080")

	os.Exit(0)
}

func Load() {

	p := flag.String("p", "127.0.0.1:80", "be proxied server address eg:  127.0.0.1:80")

	r := flag.String("r", "", "")

	dp := flag.Int("dp", 0, "remote listen port")

	debug := flag.Bool("debug", false, "debug")

	c := flag.String("c", "", "config file")

	crypt := flag.String("ct", "", "crypt type")

	k := flag.String("k", "", "crypt key")

	lp := flag.String("lp", "", "log file path")

	help := flag.Bool("help", false, "help")

	pprof := flag.String("pprof", "", "")

	flag.Parse()

	for _, v := range flag.Args() {
		if v == "help" {
			*help = true
		}
	}

	if *help {
		PrintHelp()
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
		cfg.ProxiedAddr = *p
		cfg.RemoteAddr = *r
		cfg.DestPort = *dp
		cfg.Debug = *debug
		cfg.CryptKey = *k
		cfg.LogPath = *lp
		cfg.PprofAddr = *pprof
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
