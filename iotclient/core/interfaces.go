package core

import "context"

type Signal int
type Controller interface {
	Start(ctx context.Context)
	Close() error
	SignalChannel() chan Signal
}

type Display interface {
	Display(ctx context.Context, res Response) error
	// Show(ctx context.Context, url string) error
	Close() error
	Start(ctx context.Context)
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
