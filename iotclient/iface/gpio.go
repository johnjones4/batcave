package iface

import (
	"context"
	"main/core"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stianeikeland/go-rpio"
)

type GPIOPinKey int

const (
	GPIOPinKeyToggle GPIOPinKey = iota
	GPIOPinKeyEsc
)

type GPIOConfig struct {
	SignalPins map[GPIOPinKey]int
	StatusPins map[core.StatusLight]int
}

func (c GPIOConfig) signalPins() map[GPIOPinKey]rpio.Pin {
	pins := make(map[GPIOPinKey]rpio.Pin)
	for signal, pinN := range c.SignalPins {
		pin := rpio.Pin(pinN)
		pin.Input()
		pin.PullUp()
		pins[signal] = pin
	}
	return pins
}

func (c GPIOConfig) statusPins() map[core.StatusLight]rpio.Pin {
	pins := make(map[core.StatusLight]rpio.Pin)
	for status, pinN := range c.StatusPins {
		pin := rpio.Pin(pinN)
		pin.Output()
		pin.Low()
		pins[status] = pin
	}
	return pins
}

type GPIOController struct {
	controller
	signalPins       map[GPIOPinKey]rpio.Pin
	statusPins       map[core.StatusLight]rpio.Pin
	errorOffCanceler context.CancelFunc
}

func NewGPIOController(log logrus.FieldLogger, cfg GPIOConfig) (*GPIOController, error) {
	err := rpio.Open()
	if err != nil {
		return nil, err
	}

	return &GPIOController{
		controller: controller{
			log:           log,
			signalChannel: make(chan core.Signal),
		},
		statusPins: cfg.statusPins(),
		signalPins: cfg.signalPins(),
	}, nil
}

func (c *GPIOController) Close() error {
	return rpio.Close()
}

func (c *GPIOController) Start(ctx context.Context) {
	var isTogglingDown bool
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if isTogglingDown && c.signalPins[GPIOPinKeyToggle].Read() == rpio.High {
				isTogglingDown = false
				c.signalChannel <- core.SignalTypeToggleOff
				continue
			}

			if !isTogglingDown && c.signalPins[GPIOPinKeyToggle].Read() == rpio.Low {
				isTogglingDown = true
				c.signalChannel <- core.SignalTypeToggleOn
				continue
			}

			if c.signalPins[GPIOPinKeyEsc].Read() == rpio.Low {
				c.signalChannel <- core.SignalTypeEsc
				return
			}
		}
	}
}

func (c *GPIOController) SetModeStatusLight(ctx context.Context, l core.StatusLight, t bool) error {
	switch l {
	case core.StatusLightError:
		if t {
			if c.errorOffCanceler != nil {
				c.errorOffCanceler()
			}
			ctx, c.errorOffCanceler = context.WithCancel(ctx)
			go func() {
				start := time.Now()
				for time.Since(start) < time.Second*5 {
					c.statusPins[core.StatusLightError].Toggle()
					time.Sleep(time.Millisecond * 250)
				}
				if ctx.Err() == nil {
					c.statusPins[core.StatusLightError].Low()
				}
				c.errorOffCanceler = nil
			}()
		} else {
			c.statusPins[core.StatusLightError].Low()
		}
	default:
		if t {
			c.statusPins[l].High()
		} else {
			c.statusPins[l].Low()
		}
	}
	return nil
}
