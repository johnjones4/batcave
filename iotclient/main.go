package main

import (
	"context"
	"main/iface"
	"main/util"
	"main/worker"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	controller, err := iface.NewKeyboardController(log)
	if err != nil {
		log.Panic(err)
	}

	display := iface.NewTerminalDisplayContext(log)

	cfg := util.ServerConfig{
		Hostname:        "hal9000.johnjonesfour.com",
		SecureTransport: true,
		ClientId:        "john",
		ApiKey:          "john",
	}

	rt := runtime{
		log:           log,
		controller:    controller,
		display:       display,
		lights:        display,
		voiceWorker:   worker.NewVoiceWorker(log),
		commandWorker: worker.NewCommandWorker(cfg, log),
		cfg:           cfg,
	}
	err = rt.start(context.Background())
	if err != nil {
		log.Panic(err)
	}
}
