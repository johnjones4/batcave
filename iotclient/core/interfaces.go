package core

import "context"

type Signal int
type Controller interface {
	Start(ctx context.Context)
	Close() error
	SignalChannel() chan Signal
}

type Display interface {
	Write(ctx context.Context, s string) error
	// Show(ctx context.Context, url string) error
	Close() error
}

type StatusLight int

type StatusLightsControl interface {
	SetModeStatusLight(ctx context.Context, l StatusLight, t bool) error
	Close() error
}

type Worker interface {
	Setup(errors chan error) error
	Teardown() error
}
