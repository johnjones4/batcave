package mode

import "github.com/sirupsen/logrus"

type modeConcrete struct {
	log     logrus.FieldLogger
	stop    chan bool
	stopped chan bool
}

func newModeConcrete(log logrus.FieldLogger) modeConcrete {
	return modeConcrete{
		log:     log,
		stop:    make(chan bool),
		stopped: make(chan bool),
	}
}
