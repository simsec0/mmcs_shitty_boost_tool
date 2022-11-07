package utils

import (
	"github.com/its-vichy/GoCycle"
	"log"
	"os"
	"strings"
)

var (
	Proxies *GoCycle.Cycle
	Tokens  *GoCycle.Cycle
)

func init() {
	var err error
	Proxies, err = GoCycle.NewFromFile("assets/proxies.txt")
	if err != nil {
		log.Println("[!] Can't open assets/proxies.txt!")
		os.Exit(0)
		return
	} else {
		if len(Proxies.List) == 0 {
			log.Println("[!] No proxies in assets/proxies.txt!")
			os.Exit(0)
			return
		}
	}
	Tokens, err = GoCycle.NewFromFile("assets/tokens.txt")
	if err != nil {
		log.Println("[!] Can't open assets/tokens.txt!")
		os.Exit(0)
		return
	} else {
		if len(Tokens.List) == 0 {
			log.Println("[!] No tokens in assets/tokens.txt!")
			os.Exit(0)
			return
		}
	}
	return
}

func FormatToken(token string) string {
	if strings.Contains(token, ":") {
		split := strings.Split(token, ":")
		if len(split) == 3 {
			return split[2]
		} else if len(split) == 2 {
			return split[1]
		}
	}
	return token
}
