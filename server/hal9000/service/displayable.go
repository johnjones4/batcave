package service

import "context"

type Displayable interface {
	Names() []string
	URL(context.Context) (string, error)
}

type DisplayService interface {
	Displays() []Displayable
}
