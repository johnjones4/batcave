package iface

import (
	"context"
	"main/core"

	"github.com/eiannone/keyboard"
	"github.com/sirupsen/logrus"
)

type KeyboardController struct {
	controller
}

func NewKeyboardController(log logrus.FieldLogger) (*KeyboardController, error) {
	err := keyboard.Open()
	if err != nil {
		return nil, err
	}
	return &KeyboardController{
		controller: controller{
			log:           log,
			signalChannel: make(chan core.Signal),
		},
	}, nil
}

func (c *KeyboardController) Close() error {
	return keyboard.Close()
}

func (c *KeyboardController) Start(ctx context.Context) {
	var toggleOn bool
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, key, err := keyboard.GetKey()
			if err != nil {
				c.log.Errorf("Error getting key: %s", err) //TODO
				continue
			}
			switch key {
			case keyboard.KeyEsc:
				c.signalChannel <- core.SignalTypeEsc
				return
			case keyboard.KeySpace:
				if toggleOn {
					c.signalChannel <- core.SignalTypeToggleOff
					toggleOn = false
				} else {
					c.signalChannel <- core.SignalTypeToggleOn
					toggleOn = true
				}
			}
		}
	}
}
