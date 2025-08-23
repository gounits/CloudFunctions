package tool

import (
	"github.com/BurntSushi/toml"
)

type Translate struct {
	Appid  string `toml:"appid"`
	Secret string `toml:"secret"`
	ApiKey string `toml:"apikey"`
}

type Config struct {
	Translates map[string]Translate `toml:"translates"`
}

var Conf Config

func init() {
	if _, err := toml.DecodeFile("config.toml", &Conf); err != nil {
		panic(err)
	}
}
