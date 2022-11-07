package keyauth

import (
	"github.com/Danny-Dasilva/CycleTLS/cycletls"
	"net/url"
	"strings"
)

func New(version, name, ownerId string) KeyAuth {
	client := KeyAuth{
		Client:  cycletls.Init(false),
		OwnerId: ownerId,
		Name:    name,
		Version: version,
	}
	return client
}

func (c *KeyAuth) Request(method, url, body string, headers map[string]string) (*cycletls.Response, error) {
	options := c.Options(body, headers)
	response, err := c.Client.Do(url, options, strings.ToUpper(method))
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *KeyAuth) Init() (*cycletls.Response, error) {
	headers := map[string]string{
		"content-type": "application/x-www-form-urlencoded",
	}
	params := url.Values{}
	params.Add("type", "init")
	params.Add("ver", c.Version)
	params.Add("name", c.Name)
	params.Add("ownerid", c.OwnerId)
	response, err := c.Request("POST", "https://keyauth.win/api/1.2/", params.Encode(), headers)
	if err != nil {
		return nil, err
	}
	return response, err
}

func (c *KeyAuth) Login(username, password, hwid, sessionId string) (*cycletls.Response, error) {
	headers := map[string]string{
		"content-type": "application/x-www-form-urlencoded",
	}
	params := url.Values{}
	params.Add("type", "login")
	params.Add("username", username)
	params.Add("pass", password)
	params.Add("hwid", hwid)
	params.Add("sessionid", sessionId)
	params.Add("name", c.Name)
	params.Add("ownerid", c.OwnerId)
	response, err := c.Request("POST", "https://keyauth.win/api/1.2/", params.Encode(), headers)
	if err != nil {
		return nil, err
	}
	return response, err
}

func (c *KeyAuth) Register(username, password, hwid, sessionId, licenceKey string) (*cycletls.Response, error) {
	headers := map[string]string{
		"content-type": "application/x-www-form-urlencoded",
	}
	params := url.Values{}
	params.Add("type", "register")
	params.Add("username", username)
	params.Add("pass", password)
	params.Add("key", licenceKey)
	params.Add("hwid", hwid)
	params.Add("sessionid", sessionId)
	params.Add("name", c.Name)
	params.Add("ownerid", c.OwnerId)
	response, err := c.Request("POST", "https://keyauth.win/api/1.2/", params.Encode(), headers)
	if err != nil {
		return nil, err
	}
	return response, err
}
