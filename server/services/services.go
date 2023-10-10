package services

import (
	"encoding/json"
	"main/core"
	"main/services/homeassistant"
	"main/services/llm"
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
	LLM           core.LLM
	STT           core.STT
	SocketSender  *push.SocketSender
}

type Configuration struct {
	HomeAssistantConfiguration homeassistant.HomeAssistantConfiguration `json:"homeAssistant"`
	TelegramConfiguration      telegram.TelegramConfiguration           `json:"telegram"`
	OpenAIApiKey               string                                   `json:"openAiAPIKey"`
	Ollama                     llm.OllamaCfg                            `json:"ollama"`
}

type ServiceParams struct {
	Scheduler         core.Scheduler
	Log               logrus.FieldLogger
	PushLogger        core.PushLogger
	ClientRegistry    core.ClientRegistry
	ConfigFile        string
	PushIntentFactory core.PushIntentFactory
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
		Configuration:  cfg.TelegramConfiguration,
		ClientRegistry: params.ClientRegistry,
	}

	var llmi core.LLM
	var stti core.STT
	if cfg.OpenAIApiKey != "" {
		openai := llm.NewOpenAI(params.Log, cfg.OpenAIApiKey)
		llmi = openai
		stti = openai
	}
	if cfg.Ollama.URL != "" {
		llmi = llm.NewOllama(cfg.Ollama)
	}
	// if cfg.OllamaURL != "" {
	// 	llmi = llm.NewOllama(params.Log, cfg.OllamaURL)
	// }

	socketSender := push.NewSocketSender()

	return &Services{
		TuneIn: &tunein.TuneIn{},
		HomeAssistant: &homeassistant.HomeAssistant{
			Configuration: cfg.HomeAssistantConfiguration,
			Log:           params.Log,
		},
		NOAA:     &noaa.NOAA{},
		Telegram: telegram,
		Push: &push.Push{
			ClientSenders: []core.ClientSender{
				telegram,
				socketSender,
			},
			Scheduler:         params.Scheduler,
			Log:               params.Log,
			PushLogger:        params.PushLogger,
			PushIntentFactory: params.PushIntentFactory,
			ClientRegistry:    params.ClientRegistry,
		},
		LLM:          llmi,
		STT:          stti,
		SocketSender: socketSender,
	}, nil
}
