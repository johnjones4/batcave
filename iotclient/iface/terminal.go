package iface

import (
	"context"
	"fmt"
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
	fmt.Printf("DISPLAY: \"%s\"\n", s)
	return nil
}

func (d *TerminalDisplayContext) SetModeStatusLight(ctx context.Context, l core.StatusLight, t bool) error {
	fmt.Printf("STATUS LIGHT: %s / %t\n", l.String(), t)
	return nil
}
