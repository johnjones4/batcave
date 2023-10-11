package iface

import (
	"main/core"

	"github.com/sirupsen/logrus"
)

type controller struct {
	log           logrus.FieldLogger
	signalChannel chan core.Signal
}

func (c *controller) SignalChannel() chan core.Signal {
	return c.signalChannel
}
