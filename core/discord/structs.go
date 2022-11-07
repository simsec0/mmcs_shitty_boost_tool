package discord

import (
	"discord-boosts/core/captcha"
	"github.com/Danny-Dasilva/CycleTLS/cycletls"
)

type Discord struct {
	Client            cycletls.CycleTLS
	Proxy             string
	Token             string
	Cookies           string
	Captcha           captcha.CaptchaClient
	SubscriptionSlots SubscriptionSlots
}

type SubscriptionSlots []struct {
	Id                       string      `json:"id"`
	SubscriptionId           string      `json:"subscription_id"`
	PremiumGuildSubscription interface{} `json:"premium_guild_subscription"`
	Canceled                 bool        `json:"canceled"`
	CooldownEndsAt           interface{} `json:"cooldown_ends_at"`
}
