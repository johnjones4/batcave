package service

import "github.com/johnjones4/hal-9000/server/hal9000/core"

type Displayable interface {
	Names() []string
	URL(core.Inbound) (string, error)
}

type DisplayService interface {
	Displays() []Displayable
}
