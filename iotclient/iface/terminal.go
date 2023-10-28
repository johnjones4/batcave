package iface

import (
	"context"
	"fmt"
	"main/core"

	"github.com/sirupsen/logrus"
)

type TerminalDisplay struct {
	log logrus.FieldLogger
}

func NewTerminalDisplay(log logrus.FieldLogger) *TerminalDisplay {
	return &TerminalDisplay{
		log: log,
	}
}

func (d *TerminalDisplay) Close() error {
	return nil
}

func (d *TerminalDisplay) Display(ctx context.Context, res core.Response) error {
	if res.Request != nil {
		fmt.Printf("WRITE: \"You: %s\"\n", res.Request.Message.Text)
	} else if res.PushMessage != nil {
		fmt.Printf("WRITE: \"HAL: %s\"\n", res.PushMessage.Message.Text)
		if res.PushMessage.Media.URL != "" {
			fmt.Printf("WRITE: \"HAL: %s\"\n", res.PushMessage.Media.URL)
		}
	} else if res.Response != nil {
		fmt.Printf("WRITE: \"HAL: %s\"\n", res.Response.Message.Text)
		if res.Response.Media.URL != "" {
			fmt.Printf("WRITE: \"HAL: %s\"\n", res.Response.Media.URL)
		}
	}
	return nil
}

func (d *TerminalDisplay) SetModeStatusLight(ctx context.Context, l core.StatusLight, t bool) error {
	fmt.Printf("STATUS LIGHT: %s / %t\n", l.String(), t)
	return nil
}

func (d *TerminalDisplay) Start(ctx context.Context) {
	for {
	}
}
