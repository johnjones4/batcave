package core

import "context"

type SignalType int

type ControllerSignal struct {
	Type SignalType
	Mode int
}

type Controller interface {
	Start(ctx context.Context)
	SignalChannel() chan ControllerSignal
}

type DisplayContextFactory func() (DisplayContext, error)

type DisplayContext interface {
	Write(ctx context.Context, s string) error
	Close()
	SetStatusLight(ctx context.Context, i int, t bool) error
}

type Mode interface {
	Start(ctx context.Context, displayCtx DisplayContext)
	Stop() error
	Toggle(ctx context.Context) error
}

type Worker[T any] interface {
	Start(ctx context.Context)
	Stop() error
	Chan() chan T
}
