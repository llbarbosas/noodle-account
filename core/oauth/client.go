package oauth

import "strings"

type ClientID string

// Client ...
type Client struct {
	ID                ClientID
	Label             string
	Type              string
	RedirectURI       string
	LogoutRedirectURI string
}

func (c Client) IsInternal() bool {
	return strings.Split(c.Type, ":")[0] == "internal"
}
