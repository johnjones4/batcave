package util

import (
	"net/http"
	"net/url"
)

type ServerConfig struct {
	Hostname        string
	SecureTransport bool
	ClientId        string
	ApiKey          string
}

func (c *ServerConfig) Headers() http.Header {
	return map[string][]string{
		"X-Api-Key":   {c.ApiKey},
		"X-Client-Id": {c.ClientId},
	}
}

func (c *ServerConfig) URL(path string) string {
	u := url.URL{
		Host: c.Hostname,
		Path: path,
	}
	if c.SecureTransport {
		u.Scheme = "wss"
	} else {
		u.Scheme = "ws"
	}
	return u.String()
}
