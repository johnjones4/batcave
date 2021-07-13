package main

import (
	"time"
)

type InterfaceContext interface {
	Display()
}

type InterfaceTypeSMS struct {
	PhoneNumber string
	Initiant    string
}

type InterfaceTypeTerminal struct {
	Hostname string
}

type State interface {
	TransitionToState(state State)
}

type Session struct {
	Start   time.Time
	Context InterfaceContext
	State   State
}

type Intent struct {
	Name string
}

type Action struct {
	Content string
	Session Session
	Intent  Intent
}

type Response struct {
	Message   string
	MediaURL  string
	NextState State
}
