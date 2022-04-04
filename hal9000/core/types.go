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

type RequestBody struct {
	Command  string     `json:"command"`
	Body     string     `json:"body"`
	Location Coordinate `json:"location"`
}

type Request struct {
	RequestBody
	State State `json:"state"`
}

type ResponseBody struct {
	Message string `json:"message"`
	Media   string `json:"media"`
}

type Response struct {
	ResponseBody
	State State `json:"state"`
}

type Intent interface {
	SupportedComandsForState(s State) []string
	Execute(req Request) (Response, error)
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
