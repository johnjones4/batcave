package main

import (
	"context"
	"main/api"
	"main/core"
	"main/intent"
	"main/intent/intents"
	"main/processors"
	"main/services"
	"main/store/pgstore"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
)

func main() {
	e := core.Env{}
	err := env.Parse(&e)
	if err != nil {
		panic(err)
	}

	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	store, err := pgstore.New(context.Background(), log, e.DatabaseURL)
	if err != nil {
		panic(err)
	}

	services, err := services.New(services.ServiceParams{
		Scheduler:      store,
		Log:            log,
		PushLogger:     store,
		ClientRegistry: store,
		ConfigFile:     e.ServiceConfig,
	})
	if err != nil {
		panic(err)
	}

	go services.Push.Start(context.Background())

	id, err := intent.NewIntentMatcher(log, []core.IntentActor{
		&intents.Play{
			TuneIn: services.TuneIn,
		},
		&intents.Stop{},
		&intents.ToggleDevice{
			HomeAssistant: services.HomeAssistant,
		},
		&intents.Weather{
			NOAA:      services.NOAA,
			Nominatim: services.Nominatim,
		},
		&intents.Remind{
			Push: services.Push,
		},
		&intents.Unknown{},
	}, services.LLM, store)

	if err != nil {
		panic(err)
	}

	processors := processors.Processors{
		LLM:            services.LLM,
		ClientRegistry: store,
		STT:            services.STT,
	}

	h := api.New(api.APIParams{
		IntentMatcher: id,
		RequestProcessors: []core.RequestProcessor{
			processors.SpeechToText,
			processors.DefaultLocation,
			store.LogRequest,
		},
		ResponseProcessors: []core.ResponseProcessor{
			processors.ConfirmMessage,
			store.LogResponse,
		},
		Log:            log,
		Telegram:       services.Telegram,
		ClientRegistry: store,
	})
	go h.Start(context.Background())
	err = http.ListenAndServe(e.HttpHost, h)
	panic(err)
}
