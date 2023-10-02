package core

import (
	"context"
	"net/http"
)

type RequestProcessor func(ctx context.Context, req *Request) error

type ResponseProcessor func(ctx context.Context, req *Request, res *Response) error

type ClientSender interface {
	SendToClient(ctx context.Context, clientId string, Message PushMessage) (bool, error)
}

type TelegramSender interface {
	SendOutbound(ctx context.Context, chatId int, Message OutboundMessage) error
	IsClientPermitted(ctx context.Context, r *http.Request, msgFrom int, msgText, msgType string) (bool, error)
}

type ActiveSocket struct {
	Messages     chan PushMessage
	ClientId     string
	ConnectionId string
}

type SocketSender interface {
	ClientSender
	RegisterActiveSocket(clientId string, connectionId string) *ActiveSocket
	DeregisterActiveSocket(clientId string, connectionId string) (bool, bool)
}

type Client struct {
	Source          string     `json:"source"`
	Id              string     `json:"id"`
	DefaultLocation Coordinate `json:"defaultLocation"`
	Info            any        `json:"info"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (c Coordinate) Empty() bool {
	return c.Latitude == 0 && c.Longitude == 0
}

type Message struct {
	Text  string `json:"text"`
	Audio struct {
		Data string `json:"data"`
	} `json:"audio"`
}

type Request struct {
	EventId    string     `json:"eventId"`
	Message    Message    `json:"message"`
	Source     string     `json:"source"`
	ClientID   string     `json:"clientId"`
	Coordinate Coordinate `json:"coordinate"`
}

type Media struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

type OutboundMessage struct {
	EventId string  `json:"eventId"`
	Message Message `json:"message"`
	Media   Media   `json:"media"`
}

type Response struct {
	OutboundMessage
	Action string `json:"action"`
}

type PushMessage struct {
	OutboundMessage
}
