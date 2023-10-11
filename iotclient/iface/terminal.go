package iface

import (
	"context"
	"main/core"

	"github.com/sirupsen/logrus"
)

type TerminalDisplayContext struct {
	log logrus.FieldLogger
}

func NewTerminalDisplayContext(log logrus.FieldLogger) *TerminalDisplayContext {
	return &TerminalDisplayContext{
		log: log,
	}
}

func (d *TerminalDisplayContext) Close() error {
	return nil
}

func (d *TerminalDisplayContext) Write(ctx context.Context, s string) error {
	d.log.Printf("DISPLAY: \"%s\"", s)
	return nil
}

func (d *TerminalDisplayContext) SetModeStatusLight(ctx context.Context, l core.StatusLight, t bool) error {
	d.log.Printf("STATUS LIGHT: %d / %b", l, t)
	return nil
}
