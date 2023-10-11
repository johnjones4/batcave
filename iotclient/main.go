package main

import (
	"context"
	"flag"
	"main/core"
	"main/iface"
	"main/util"
	"main/worker"

	"github.com/sirupsen/logrus"
)

const (
	modeTerminal  = "terminal"
	modeInterface = "interface"
)

func main() {
	log := logrus.New()

	mode := flag.String("mode", modeTerminal, "Operating mode")
	loglevelStr := flag.String("loglevel", logrus.WarnLevel.String(), "Log level")
	flag.Parse()

	loglevel, err := logrus.ParseLevel(*loglevelStr)
	if err != nil {
		log.Panic(err)
	}
	log.SetLevel(loglevel)

	var controller core.Controller
	var display core.Display
	var lights core.StatusLightsControl

	switch *mode {
	case modeTerminal:
		controller, err = iface.NewKeyboardController(log)
		if err != nil {
			log.Panic(err)
		}
		term := iface.NewTerminalDisplayContext(log)
		display = term
		lights = term
	case modeInterface:
		gpio, err := iface.NewGPIOController(log, iface.GPIOConfig{
			StatusPins: map[core.StatusLight]int{},
			SignalPins: map[iface.GPIOPinKey]int{},
		})
		if err != nil {
			log.Panic(err)
		}
		controller = gpio
		lights = gpio
		display = iface.NewTerminalDisplayContext(log)
	default:
		return
	}

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
		lights:        lights,
		voiceWorker:   worker.NewVoiceWorker(log),
		commandWorker: worker.NewCommandWorker(cfg, log),
		cfg:           cfg,
	}
	err = rt.start(context.Background())
	if err != nil {
		log.Panic(err)
	}
}
