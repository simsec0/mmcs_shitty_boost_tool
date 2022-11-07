package main

import (
	"bufio"
	"crypto/md5"
	"discord-boosts/core/console"
	"discord-boosts/core/discord"
	"discord-boosts/core/keyauth"
	"discord-boosts/core/tasks"
	"discord-boosts/core/utils"
	"encoding/hex"
	"fmt"
	"github.com/howeyc/gopass"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"
)

var (
	Version   string = "1.0"
	Name      string = "boost tool"
	OwnerId   string = "NE4QikZYZu"
	HWID      string
	SessionID string

	username   string = "guest"
	authorized bool   = false

	clients    []discord.Discord
	initalized bool = false

	commands = []string{
		"exit",
		"help",
		"login",
		"register",
		"reset",
		"boost",
		"stats",
		"setup",
	}

	authorizationCommands = []string{
		"help",
		"boost",
		"setup",
	}
)

func init() {
	hasher := md5.New()
	name, _ := os.Hostname()
	usr, _ := user.Current()
	hasher.Write([]byte(name + usr.Username))
	HWID = hex.EncodeToString(hasher.Sum(nil))
}

func canUse(command string) bool {
	for _, content := range authorizationCommands {
		if command == content {
			if authorized {
				return true
			} else {
				return false
			}
		}
	}

	return true
}

//
//func start(guildId, invite string) {
//	var wg sync.WaitGroup
//
//	conc := make(chan struct{}, 100)
//
//	for l := 0; l < len(utils.Tokens.List); l++ {
//		wg.Add(1)
//
//		go func(token string) {
//			defer wg.Done()
//			conc <- struct{}{}
//
//			cookies, err := discord.GetCookies(fmt.Sprintf("http://%v", utils.Proxies.Next()))
//			if err != nil {
//				log.Println("[!] Failed to fetch cookies.")
//				return
//			}
//			client := discord.New(token, cookies, utils.Proxies.Next(), true)
//			client.Join(invite)
//			response, err := client.CheckSubscriptionSlots()
//			if err != nil {
//				log.Println("[!] Failed to check subscriptions.")
//				return
//			}
//
//			boosts := []map[string]any{}
//
//			err = json2.Unmarshal([]byte(response.Body), &boosts)
//			if err != nil {
//				log.Println(err)
//				log.Println("[!] Failed to unmarshal subscriptions.")
//				return
//			}
//			if len(boosts) != 0 {
//				for i := 0; i < len(boosts); i++ {
//					resp, err := client.BoostServer(guildId, boosts[i]["id"].(string))
//					if err != nil {
//						log.Printf("[!] Failed to boost: %v", err)
//						continue
//					}
//					log.Printf("[!] Successfully boosted %v with user %v", resp.JSONBody()["guild_id"], resp.JSONBody()["user_id"])
//				}
//			} else {
//				log.Println("[!] No boosts available!")
//			}
//			<-conc
//		}(utils.Tokens.Next())
//
//	}
//	//
//	wg.Wait()
//}

func Clear() {
	cmd := exec.Command("cmd", "/c", "cls||clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	return
}

func main() {
	console.Clear()
	key := keyauth.New(Version, Name, OwnerId)
	if !authorized {
		console.SetTitle("Boost - Please login to your account!")
		resp, _ := key.Init()
		if resp.Body == "" {
			console.Print("\u001B[0mKeyAuth so shit. Refreshing... \n")
			time.Sleep(2 * time.Second)
			main()
		}
		SessionID = resp.JSONBody()["sessionid"].(string)
	} else {
		console.SetTitle(fmt.Sprintf("Boost - Logged in as %v!", username))
	}

	console.ShowBanner()
	if !authorized {
		fmt.Printf("\u001B[0m%v Hey! You aren't logged in, please use \033[1m\033[4m\u001B[38;5;59mlogin\u001B[0m or \u001B[1m\u001B[4m\u001B[38;5;59mregister\u001B[0m.\n\n", strings.Repeat(" ", 30))
	} else {
		fmt.Printf("\u001B[0m%v Hey %v! To see all the available commands try using \033[1m\033[4m\u001B[38;5;59mhelp\u001B[0m.\n\n", strings.Repeat(" ", 25), username)
	}

	for {
		command := console.PromptInput(username)

		if !canUse(command) {
			console.Print("You need to be logged in to use this command.\n")
			continue
		}

		if command == "help" {
			console.Print("\u001B[1mAvailable commands:\u001B[0m\n")
			for _, name := range commands {
				console.Print(fmt.Sprintf("%v\u001B[1m\u001B[38;5;59m-\u001B[0m %v\n", strings.Repeat(" ", 3), name))
			}
			fmt.Println()
		} else if command == "reset" {
			main()
			break
		} else if command == "exit" {
			os.Exit(0)
		} else if command == "login" {
			if authorized {
				console.Print("You are already logged in.\n")
				continue
			}

			console.Print("\u001B[0mUsername: ")
			reader := bufio.NewReader(os.Stdin)
			inputUsername, err := reader.ReadString('\n')
			if err != nil {
				main()
				continue
			}

			inputUsername = strings.TrimSuffix(inputUsername, "\n")

			console.Print("\u001B[0mPassword: ")
			password, _ := gopass.GetPasswdMasked()

			response, err := key.Login(inputUsername, string(password), HWID, SessionID)

			if err != nil {
				main()
				break
			}

			if response.Body == "" {
				console.Print("\u001B[0mKeyAuth so shit. Refreshing... \n")
				time.Sleep(2 * time.Second)
				main()
				break
			}
			if response.JSONBody()["success"] != true {
				console.Print("\u001B[0mInvalid Details Provided.\n")
				time.Sleep(5 * time.Second)
				main()
				break
			} else {
				console.Print(fmt.Sprintf("\u001B[0mHey %v! Welcome back to boosts!\n", response.JSONBody()["info"].(map[string]interface{})["username"]))
				authorized = true
				time.Sleep(2 * time.Second)
				username = "root"
				main()
				break
			}
		} else if command == "register" {
			if authorized {
				console.Print("You are already logged in.\n")
				continue
			}
			console.Print("\u001B[0mUsername: ")
			reader := bufio.NewReader(os.Stdin)
			inputUsername, err := reader.ReadString('\n')
			if err != nil {
				main()
				continue
			}
			inputUsername = strings.TrimSuffix(inputUsername, "\n")
			console.Print("\u001B[0mPassword: ")
			password, _ := gopass.GetPasswdMasked()
			console.Print("\u001B[0mLicense Key: ")
			reader = bufio.NewReader(os.Stdin)
			licenceKey, err := reader.ReadString('\n')
			if err != nil {
				main()
				continue
			}
			resp, err := key.Register(inputUsername, string(password), HWID, SessionID, licenceKey)
			if err != nil {
				main()
				continue
			}
			if resp.Body == "" {
				console.Print("\u001B[0mKeyAuth so shit. Refreshing... \n")
				time.Sleep(2 * time.Second)
				main()
				break
			}

			if resp.JSONBody()["success"] == true {
				console.Print("\u001B[0mWelcome to boosts.\n")
				time.Sleep(1 * time.Second)
				authorized = true
				username = "root"
				main()
				break
			} else {
				console.Print(fmt.Sprintf("\u001B[0m%v\n", resp.JSONBody()["message"]))
				time.Sleep(1 * time.Second)
				main()
				break
			}

		} else if command == "stats" {
			console.Print(fmt.Sprintf("You currently have \u001B[1m\u001B[4m\u001B[38;5;59m%v\u001B[0m tokens.\n", len(utils.Tokens.List)))
			console.Print(fmt.Sprintf("You currently have \u001B[1m\u001B[4m\u001B[38;5;59m%v\u001B[0m proxies.\n", len(utils.Proxies.List)))
			fmt.Println()
		} else if command == "boost" {
			if !initalized || len(clients) == 0 {
				console.Print("Please run \u001B[38;5;59msetup\u001B[0m before executing this command!\n")
				continue
			}
			var guildId string
			var invite string
			var amount int

			console.Print("Guild ID: ")
			fmt.Scanln(&guildId)

			console.Print("Server Invite: discord.gg/")
			fmt.Scanln(&invite)

			console.Print("Amount of boosts: ")
			fmt.Scanln(&amount)

			start := time.Now()
			tasks.Join(invite, clients, amount/2)
			tasks.Boost(guildId, clients, amount/2)
			fmt.Println()
			console.Print(fmt.Sprintf("Finished task in %v.\n", time.Since(start)))
		} else if command == "setup" {
			if initalized {
				console.Print("You have already setup all the clients.\n")
				continue
			}
			fmt.Println()

			clients = tasks.Setup()
			initalized = true
			fmt.Println()
			console.Print(fmt.Sprintf("Successfully setup \u001B[38;5;59m%v\u001B[0m clients.\n", len(clients)))
		}
	}
}
