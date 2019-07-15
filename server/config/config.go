package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
)

var Port int
var Debug bool

type config struct {
	Port  int  `json:"port"`
	Debug bool `json:"debug"`
}

func Load() {

	client_port := flag.Int("client_port", 0, "")
	debug := flag.Bool("debug", false, "")
	c := flag.String("c", "", "config file")

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

		Debug = cfg.Debug
		Port = cfg.Port

	} else {

		Port = *client_port
		Debug = *debug
	}

}
