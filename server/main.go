package main

import (
	"context"
	"main/api"
	"main/core"
	"main/intent"
	"main/intent/intents"
	"main/llm"
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

	store, err := pgstore.New(context.Background(), e.DatabaseURL)
	if err != nil {
		panic(err)
	}

	// llm := llm.NewOllama("http://localhost:11434", log)
	llm := llm.NewOpenAI(log, e.OpenAIKey)

	if e.IntentDescriptions != "" {
		err = intent.ParseIntents(context.Background(), log, llm, store, e.IntentDescriptions)
		if err != nil {
			panic(err)
		}
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
	}, llm, store)

	if err != nil {
		panic(err)
	}

	processors := processors.Processors{
		LLM: llm,
	}

	h := api.New(api.APIParams{
		IntentMatcher: id,
		RequestProcessors: []core.RequestProcessor{
			store.LogRequest,
		},
		ResponseProcessors: []core.ResponseProcessor{
			store.LogResponse,
			processors.ConfirmMessage,
		},
		Log:      log,
		Telegram: *services.Telegram,
	})
	err = http.ListenAndServe(e.HttpHost, h)
	panic(err)
}
