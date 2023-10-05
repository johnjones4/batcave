package core

import "errors"

const NStatusLights = 3

var (
	ErrorDisplayContextClosed = errors.New("display context closed")
)

const (
	SignalTypeMode SignalType = iota
	SignalTypeToggle
	SignalTypeEsc
)
