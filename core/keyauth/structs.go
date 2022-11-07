package keyauth

import "github.com/Danny-Dasilva/CycleTLS/cycletls"

type KeyAuth struct {
	Client  cycletls.CycleTLS
	OwnerId string
	Version string
	Name    string
}
