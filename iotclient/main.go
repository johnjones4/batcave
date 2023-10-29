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
	inputModeKeyboard  = "keyboard"
	inputModeGPIO      = "gpio"
	outputModeTerminal = "terminal"
	outputModeGUI      = "gui"
	outputModeGPIO     = "gpio"
)

func main() {
	log := logrus.New()

	inputMode := flag.String("inputmode", inputModeKeyboard, "Input mode")
	outputMode := flag.String("outputmode", outputModeTerminal, "Output mode")
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
	var gpio *iface.GPIOController

	switch *inputMode {
	case inputModeKeyboard:
		controller, err = iface.NewKeyboardController(log)
		if err != nil {
			log.Panic(err)
		}
	case inputModeGPIO:
		gpio, err = iface.NewGPIOController(log, iface.GPIOConfig{
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
		controller = gpio
		if err != nil {
			log.Panic(err)
		}
	default:
		return
	}
	switch *outputMode {
	case outputModeTerminal:
		term := iface.NewTerminalDisplay(log)
		display = term
		lights = term
	case outputModeGUI:
		lights = gpio
		display = iface.NewGUIDisplay()
	case outputModeGPIO:
		term := iface.NewTerminalDisplay(log)
		display = term
		lights = gpio
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
		splayer:       worker.NewStreamPlayer(log),
		bplayer:       worker.NewBufferPlayer(log),
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
