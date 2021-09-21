package main

const (
	ConnectionTypeSocket    = "socket"
	ConnectionTypeWebsocket = "websocket"
)

type ConnectionEvent struct {
	Source string `json:"source"`
	Type   string `json:"type"`
}
