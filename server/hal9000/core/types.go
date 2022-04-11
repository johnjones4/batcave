package core

import "time"

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type State struct {
	User  User   `json:"user"`
	State string `json:"state"`
}

type InboundBody struct {
	Body     string     `json:"body"`
	Location Coordinate `json:"location"`
}

type Inbound struct {
	InboundBody
	Command string `json:"command"`
	State   State  `json:"state"`
}

type OutboundBody struct {
	Body  string `json:"body"`
	Media string `json:"media"`
	URL   string `json:"url"`
}

type Outbound struct {
	OutboundBody
	State State `json:"state"`
}

type CommandInfo struct {
	Description  string `json:"description"`
	RequiresBody bool   `json:"requiresBody"`
}

type Intent interface {
	SupportedComandsForState(s State) map[string]CommandInfo
	Execute(req Inbound) (Outbound, error)
}

type FeedbackError struct {
	message string
}

type User struct {
	Name string `json:"name"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token      string    `json:"token"`
	User       string    `json:"user"`
	Expiration time.Time `json:"expiration"`
}
