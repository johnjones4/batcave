package services

import (
	"encoding/json"
	"main/core"
	"main/services/homeassistant"
	"main/services/noaa"
	"main/services/nominatim"
	"main/services/push"
	"main/services/telegram"
	"main/services/tunein"
	"os"

	"github.com/sirupsen/logrus"
)

type Services struct {
	HomeAssistant *homeassistant.HomeAssistant
	TuneIn        *tunein.TuneIn
	NOAA          *noaa.NOAA
	Nominatim     *nominatim.Nominatim
	Telegram      *telegram.Telegram
	Push          *push.Push
}

type Configuration struct {
	HomeAssistantConfiguration homeassistant.HomeAssistantConfiguration `json:"homeAssistant"`
	TelegramToken              string                                   `json:"telegramToken"`
	PushConfiguration          push.PushConfiguration                   `json:"push"`
}

type ServiceParams struct {
	Scheduler      core.Scheduler
	Log            logrus.FieldLogger
	PushLogger     core.PushLogger
	ClientRegistry core.ClientRegistry
	ConfigFile     string
}

func New(params ServiceParams) (*Services, error) {
	bytes, err := os.ReadFile(params.ConfigFile)
	if err != nil {
		return nil, err
	}

	var cfg Configuration
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	telegram := &telegram.Telegram{
		Token:          cfg.TelegramToken,
		ClientRegistry: params.ClientRegistry,
	}

	return &Services{
		TuneIn: &tunein.TuneIn{},
		HomeAssistant: &homeassistant.HomeAssistant{
			Configuration: cfg.HomeAssistantConfiguration,
			Log:           params.Log,
		},
		NOAA:     &noaa.NOAA{},
		Telegram: telegram,
		Push: &push.Push{
			PushConfiguration: cfg.PushConfiguration,
			ClientProviders: []core.ClientProvider{
				telegram,
			},
			Scheduler:  params.Scheduler,
			Log:        params.Log,
			PushLogger: params.PushLogger,
		},
	}, nil
}
