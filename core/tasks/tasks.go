package tasks

import (
	"discord-boosts/core/assets"
	"discord-boosts/core/captcha"
	"discord-boosts/core/console"
	"discord-boosts/core/discord"
	"discord-boosts/core/utils"
	"fmt"
	"github.com/Danny-Dasilva/CycleTLS/cycletls"
	"github.com/zenthangplus/goccm"
)

func Setup() []discord.Discord {
	concs := goccm.New(500)
	var clients []discord.Discord

	for i := 1; i <= len(utils.Tokens.List)+1; i++ {
		concs.Wait()

		go func() {
			client := discord.Discord{
				Client:  cycletls.Init(false),
				Captcha: captcha.CaptchaClient{Apikey: assets.Configuration.CaptchaKey, Host: "https://api.capmonster.cloud"},
				Token:   utils.FormatToken(utils.Tokens.Next()),
				Proxy:   utils.Proxies.Next(),
			}
			cookies, err := client.GetCookies()
			if err != nil {
				console.Print(fmt.Sprintf("\u001B[38;5;59m[!]\u001B[0m An error has occured \u001B[38;5;59m->\u001B[0m %v\n", err.Error()))
			} else if cookies == "" {
				console.Print("\u001B[38;5;59m[!]\u001B[0m Failed to get cookies.\n")
			} else {
				console.Print("\u001B[38;5;59m[!]\u001B[0m Fetched cookies.\n")
				clients = append(clients, client)
			}
			concs.Done()
		}()
	}
	concs.WaitAllDone()
	return clients
}

func Join(invite string, clients []discord.Discord, amount ...int) {
	conc := goccm.New(500)
	if len(amount) != 0 {
		for i := 1; i <= amount[0]; i++ {
			conc.Wait()
			go func(client discord.Discord) {
				client.Join(invite)
				conc.Done()
			}(clients[i-1])
		}
		conc.WaitAllDone()
	} else {
		for i := 1; i <= len(clients); i++ {
			conc.Wait()
			go func(client discord.Discord) {
				client.Join(invite)
				conc.Done()
			}(clients[i-1])
		}
		conc.WaitAllDone()
	}
}

func Boost(guildId string, clients []discord.Discord, amount ...int) {
	conc := goccm.New(500)
	if len(amount) != 0 {
		for i := 1; i <= amount[0]; i++ {
			conc.Wait()
			go func(client discord.Discord) {
				client.CheckSubscriptionSlots()
				for _, slot := range client.SubscriptionSlots {
					_, err := client.BoostServer(guildId, slot.Id)
					if err != nil {
						console.Print(fmt.Sprintf("\u001B[38;5;59m[!]\u001B[0m Failed to boost. %v\n", err))
						return
					}
					console.Print("\u001B[38;5;59m[!]\u001B[0m Successfully boosted.\n")
				}
			}(clients[i-1])
		}
		conc.WaitAllDone()
	}
}
