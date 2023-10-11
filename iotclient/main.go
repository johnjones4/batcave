package main

import (
	"context"
	"flag"
	"main/core"
	"main/iface"
	"main/util"
	"main/worker"
	"os"
	"strings"

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
			StatusPins: map[core.StatusLight]int{
				core.StatusLightError:     17,
				core.StatusLightListening: 27,
				core.StatusLightWorking:   22,
			},
			SignalPins: map[iface.GPIOPinKey]int{
				iface.GPIOPinKeyEsc:    5,
				iface.GPIOPinKeyToggle: 6,
			},
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
		Hostname:        os.Getenv("HOSTNAME"),
		SecureTransport: parseBool(os.Getenv("SECURE_TRANSPORT")),
		ClientId:        os.Getenv("CLIENT_ID"),
		ApiKey:          os.Getenv("API_KEY"),
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

func parseBool(b string) bool {
	return strings.ToLower(b) == "true" || b == "1"
}
