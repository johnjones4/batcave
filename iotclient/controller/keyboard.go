package controller

import (
	"context"
	"main/core"

	"github.com/eiannone/keyboard"
	"github.com/sirupsen/logrus"
)

type KeyboardController struct {
	log           logrus.FieldLogger
	signalChannel chan core.ControllerSignal
}

func NewKeyboardController(log logrus.FieldLogger) *KeyboardController {
	return &KeyboardController{
		log:           log,
		signalChannel: make(chan core.ControllerSignal),
	}
}

func (c *KeyboardController) SignalChannel() chan core.ControllerSignal {
	return c.signalChannel
}

func (c *KeyboardController) Start(ctx context.Context) {
	err := keyboard.Open()
	if err != nil {
		c.log.Errorf("Error getting keyboard: %s", err)
		return
	}
	defer func() {
		keyboard.Close()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			char, key, err := keyboard.GetKey()
			if err != nil {
				c.log.Errorf("Error getting key: %s", err)
				continue
			}
			switch key {
			case keyboard.KeyEsc:
				c.signalChannel <- core.ControllerSignal{
					Type: core.SignalTypeEsc,
				}
				return
			}
			c.log.Debugf("Got key %c", char)
			switch char {
			case rune('1'):
				c.signalChannel <- core.ControllerSignal{
					Type: core.SignalTypeMode,
					Mode: 0,
				}
			case rune('2'):
				c.signalChannel <- core.ControllerSignal{
					Type: core.SignalTypeMode,
					Mode: 1,
				}
			case rune('p'):
				c.signalChannel <- core.ControllerSignal{
					Type: core.SignalTypeToggle,
				}
			}
		}
	}
}
