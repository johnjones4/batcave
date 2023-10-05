package main

import (
	"context"
	"main/controller"
	"main/core"
	display "main/displaycontext"
	"main/mode"
	"main/util"
	"main/worker"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	controller := controller.NewKeyboardController(log)
	cfg := util.ServerConfig{
		Hostname:        "hal9000.johnjonesfour.com",
		SecureTransport: true,
		ClientId:        "john",
		ApiKey:          "john",
	}
	logWorker := worker.NewLogWorker(cfg, log)
	commandWorker := worker.NewCommandWorker(cfg, log)
	voiceWorker := worker.NewVoiceWorker(log)
	err := voiceWorker.Setup()
	if err != nil {
		log.Panic(err)
	}
	modes := []core.Mode{
		mode.NewModeLog(log, logWorker),
		mode.NewModeCommand(log, voiceWorker, commandWorker, cfg.ClientId),
	}
	rt := New(log, modes, controller, func() (core.DisplayContext, error) {
		return display.NewTerminalDisplayContext(log), nil
	})
	rt.Start(context.Background())
}
