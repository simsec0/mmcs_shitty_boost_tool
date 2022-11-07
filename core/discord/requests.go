package discord

import (
	"fmt"
	"github.com/Danny-Dasilva/CycleTLS/cycletls"
)

func (c *Discord) Options(body string, headers map[string]string) cycletls.Options {
	options := cycletls.Options{
		Body:      body,
		Ja3:       "771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-17513,29-23-24,0",
		Headers:   headers,
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) discord/1.0.9006 Chrome/91.0.4472.164 Electron/13.6.6 Safari/537.36",
		Proxy:     fmt.Sprintf("http://%v", c.Proxy),
	}
	return options
}
