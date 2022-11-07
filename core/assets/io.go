package assets

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

func LoadConfig() *Config {
	var reply Config
	_, err := toml.DecodeFile("assets/config.toml", &reply)
	if err != nil {
		log.Println("[!] Failed to open assets/config.toml!")
		os.Exit(0)
		return nil
	}
	return &reply
}
