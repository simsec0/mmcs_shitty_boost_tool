package discord

import (
	"discord-boosts/core/assets"
	"discord-boosts/core/captcha"
	"discord-boosts/core/console"
	"encoding/json"
	"fmt"
	"github.com/Danny-Dasilva/CycleTLS/cycletls"
	"strings"
)

func New(token, cookies, proxy string, workers bool) Discord {
	client := Discord{
		Client:  cycletls.Init(workers),
		Proxy:   proxy,
		Token:   token,
		Cookies: cookies,
		Captcha: captcha.CaptchaClient{Apikey: assets.Configuration.CaptchaKey, Host: "https://api.capmonster.cloud"},
	}
	return client
}

func (c *Discord) Request(method, url, body string, headers map[string]string) (*cycletls.Response, error) {
	options := c.Options(body, headers)
	response, err := c.Client.Do(url, options, strings.ToUpper(method))
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Discord) Join(invite string, captchaKey ...string) {
	var payload string
	headers := map[string]string{
		"accept":             "*/*",
		"accept-encoding":    "gzip, deflate, br",
		"accept-language":    "en-GB,en-NZ;q=0.9",
		"authorization":      c.Token,
		"cookie":             c.Cookies,
		"content-type":       "application/json",
		"origin":             "https://discord.com",
		"referer":            "https://google.com",
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-origin",
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) discord/1.0.9006 Chrome/91.0.4472.164 Electron/13.6.6 Safari/537.36",
		"x-debug-options":    "bugReporterEnabled",
		"x-discord-locale":   "en-US",
		"x-super-properties": "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRGlzY29yZCBDbGllbnQiLCJyZWxlYXNlX2NoYW5uZWwiOiJzdGFibGUiLCJjbGllbnRfdmVyc2lvbiI6IjEuMC45MDA2Iiwib3NfdmVyc2lvbiI6IjEwLjAuMTkwNDQiLCJvc19hcmNoIjoieDY0Iiwic3lzdGVtX2xvY2FsZSI6ImVuLUdCIiwiY2xpZW50X2J1aWxkX251bWJlciI6MTQxNDcxLCJjbGllbnRfZXZlbnRfc291cmNlIjpudWxsfQ==",
	}

	if len(captchaKey) == 0 {
		payload = "{}"
	} else {
		payload = fmt.Sprintf(`{"captcha_key": "%v", "captcha_rqtoken": "%v"}`, captchaKey[0], captchaKey[1])
	}

	response, err := c.Request("POST", fmt.Sprintf("https://discord.com/api/v9/invites/%v", invite), payload, headers)
	if err != nil {
		console.Print("\u001B[38;5;59m[!]\u001B[0m Token failed to join server, retrying...\n")
		c.Join(invite)
		return
	}
	if strings.Contains(response.Body, "new_member") {
		console.Print("\u001B[38;5;59m[!]\u001B[0m Successfully joined server!\n")
		return
	} else if strings.Contains(response.Body, "new_member") {
		console.Print("\u001B[38;5;59m[!]\u001B[0m Already in server!\n")
		return
	} else if strings.Contains(response.Body, "captcha") {
		console.Print("\u001B[38;5;59m[!]\u001B[0m Captcha required. Solving...\n")

		jsonResponse := response.JSONBody()

		task := captcha.Task{
			Type:    "HCaptchaTaskProxyless",
			Host:    "https://discord.com/",
			Sitekey: jsonResponse["captcha_sitekey"].(string),

			Data:      jsonResponse["captcha_rqdata"].(string),
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) discord/1.0.9006 Chrome/91.0.4472.164 Electron/13.6.6 Safari/537.36",
		}

		taskId, err := c.Captcha.CreateTask(task)
		if err != nil {
			console.Print("\u001B[38;5;59m[!]\u001B[0m Failed to create captcha task. Retrying...\n")
			c.Join(invite)
			return
		}

		console.Print(fmt.Sprintf("\u001B[38;5;59m[!]\u001B[0m Captcha task created! Task ID: %v\n", taskId.TaskID))
		captchaKey, err := c.Captcha.JoinTaskResult(taskId.TaskID)
		if err != nil {
			console.Print("\u001B[38;5;59m[!]\u001B[0m Token failed to join server, retrying...\n")
			c.Join(invite)
			return
		}
		console.Print(fmt.Sprintf("\u001B[38;5;59m[!]\u001B[0m Captcha solved! %v\n", taskId.TaskID))
		c.Join(invite, captchaKey.Solution.GeneratedPassUUID, jsonResponse["captcha_rqtoken"].(string))
		return
	} else {
		console.Print("\u001B[38;5;59m[!]\u001B[0m Token failed to join server.\n")
	}

	if response.Status == 200 {
		console.Print("\u001B[38;5;59m[!]\u001B[0m Successfully joined server.\n")
		return
	} else {
		console.Print("\u001B[38;5;59m[!]\u001B[0m Failed to join server.\n")
		return
	}
}

func (c *Discord) CheckSubscriptionSlots() {
	headers := map[string]string{
		"accept":           "*/*",
		"accept-encoding":  "gzip, deflate, br",
		"accept-language":  "en-GB,en-NZ;q=0.9",
		"authorization":    c.Token,
		"cookie":           c.Cookies,
		"content-type":     "application/json",
		"origin":           "https://discord.com",
		"referer":          "https://discord.com",
		"sec-fetch-dest":   "empty",
		"sec-fetch-mode":   "cors",
		"sec-fetch-site":   "same-origin",
		"user-agent":       "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) discord/1.0.9006 Chrome/91.0.4472.164 Electron/13.6.6 Safari/537.36",
		"x-debug-options":  "bugReporterEnabled",
		"x-discord-locale": "en-US",
	}

	response, err := c.Request("GET", "https://discord.com/api/v9/users/@me/guilds/premium/subscription-slots", "", headers)
	if err != nil {
		console.Print(fmt.Sprintf("\u001B[38;5;59m[!]\u001B[0m %v\n", err))
		return
	}

	var slots SubscriptionSlots
	err = json.Unmarshal([]byte(response.Body), &slots)
	if err != nil {
		console.Print(fmt.Sprintf("\u001B[38;5;59m[!]\u001B[0m %v\n", err))
		return
	}

	c.SubscriptionSlots = slots
	return

}

func (c *Discord) BoostServer(guildID, slotId string) (*cycletls.Response, error) {
	var payload = fmt.Sprintf(`{"user_premium_guild_subscription_slot_ids": ["%v"]}`, slotId)
	headers := map[string]string{
		"accept":           "*/*",
		"accept-encoding":  "gzip, deflate, br",
		"accept-language":  "en-GB,en-NZ;q=0.9",
		"authorization":    c.Token,
		"cookie":           c.Cookies,
		"content-type":     "application/json",
		"origin":           "https://discord.com",
		"referer":          "https://discord.com",
		"sec-fetch-dest":   "empty",
		"sec-fetch-mode":   "cors",
		"sec-fetch-site":   "same-origin",
		"user-agent":       "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) discord/1.0.9006 Chrome/91.0.4472.164 Electron/13.6.6 Safari/537.36",
		"x-debug-options":  "bugReporterEnabled",
		"x-discord-locale": "en-US",
	}
	response, err := c.Request("PUT", fmt.Sprintf("https://discord.com/api/v9/guilds/%v/premium/subscriptions", guildID), payload, headers)
	if err != nil {
		return nil, err
	}
	return response, err
}

func (c *Discord) GetCookies() (string, error) {
	headers := map[string]string{
		"accept":          "*/*",
		"accept-encoding": "gzip, deflate, br",
		"accept-language": "en-GB,en-NZ;q=0.9",
		"content-type":    "application/json",
		"origin":          "https://discord.com",
		"referer":         "https://google.com",
		"sec-fetch-dest":  "empty",
		"sec-fetch-mode":  "cors",
		"sec-fetch-site":  "same-origin",
		"user-agent":      "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) discord/1.0.9006 Chrome/91.0.4472.164 Electron/13.6.6 Safari/537.36",
	}
	response, err := c.Request("GET", "https://discord.com/", "", headers)
	if err != nil {
		return "", err
	}
	cookies := response.Headers["Set-Cookie"]
	c.Cookies = cookies
	return cookies, nil
}
