package display

import (
	"context"
	"main/core"

	"github.com/sirupsen/logrus"
)

type TerminalDisplayContext struct {
	log    logrus.FieldLogger
	closed bool
}

func NewTerminalDisplayContext(log logrus.FieldLogger) *TerminalDisplayContext {
	return &TerminalDisplayContext{
		log: log,
	}
}

func (d *TerminalDisplayContext) Close() {
	d.log.Debug("Closing display")
	d.closed = true
}

func (d *TerminalDisplayContext) Write(ctx context.Context, s string) error {
	if d.closed {
		return core.ErrorDisplayContextClosed
	}
	d.log.Printf("DISPLAY: \"%s\"", s)
	return nil
}

func (d *TerminalDisplayContext) SetStatusLight(ctx context.Context, i int, t bool) error {
	if d.closed {
		return core.ErrorDisplayContextClosed
	}
	if i < core.NStatusLights {
		d.log.Printf("DISPLAY: Status light %d is on: %b", i, t)
	}
	return nil
}
