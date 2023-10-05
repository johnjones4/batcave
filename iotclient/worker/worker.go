package worker

import "github.com/sirupsen/logrus"

type workerConcrete struct {
	log     logrus.FieldLogger
	stop    chan bool
	stopped chan bool
}

func newWorkerConcrete(log logrus.FieldLogger) workerConcrete {
	return workerConcrete{
		log:     log,
		stop:    make(chan bool),
		stopped: make(chan bool),
	}
}
