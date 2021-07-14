package main

import (
	"time"
)

type InterfaceContext interface {
	Display()
}

type InterfaceTypeSMS struct {
	Person   Person
	Initiant string
}

type InterfaceTypeTerminal struct {
	Machine Machine
}

type Response struct {
	Message  string
	MediaURL string
}

type Noun interface {
	Act(object interface{}) (State, Response, error)
}

type Person struct {
	PhoneNumber string
	Email       string
}

type Machine struct {
}

type ParsedInputMessage struct {
	Message  string
	Segments []interface{}
}

type Intent interface {
	Name() string
	Process(p ParsedInputMessage) (State, Response, error)
}

type State interface {
	Name() string
	CanTransitionToState(state State) bool
	InferIntent(p ParsedInputMessage) (Intent, error)
}

type Session struct {
	Start   time.Time
	Context InterfaceContext
	State   State
}
