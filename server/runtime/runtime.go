package runtime

import (
	"context"
	"main/api"
	"main/core"
	"main/intent"
	"main/processors"
	"main/services"
	"main/store/pgstore"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
)

type Runtime struct {
	Env      core.Env
	Log      *logrus.Logger
	Store    *pgstore.PGStore
	Services *services.Services
	API      *api.API
	Intents  []core.IntentActor
}

func New(ctx context.Context) (*Runtime, error) {
	r := &Runtime{
		Log: logrus.New(),
	}
	r.Log.SetLevel(logrus.DebugLevel)

	err := env.Parse(&r.Env)
	if err != nil {
		return nil, err
	}

	r.Store, err = pgstore.New(ctx, r.Log, r.Env.DatabaseURL)
	if err != nil {
		return nil, err
	}

	r.Services, err = services.New(services.ServiceParams{
		Scheduler:         r.Store,
		Log:               r.Log,
		PushLogger:        r.Store,
		ClientRegistry:    r.Store,
		ConfigFile:        r.Env.ServiceConfig,
		PushIntentFactory: r,
	})
	if err != nil {
		return nil, err
	}

	r.Intents = intent.Factory(r.Services)
	intentMatcher, err := intent.NewIntentMatcher(r.Log, r.Intents, r.Services.LLM, r.Store)
	if err != nil {
		return nil, err
	}

	processors := processors.Processors{
		LLM:            r.Services.LLM,
		ClientRegistry: r.Store,
		STT:            r.Services.STT,
	}

	r.API = api.New(api.APIParams{
		IntentMatcher: intentMatcher,
		RequestProcessors: []core.RequestProcessor{
			processors.SpeechToText,
			processors.DefaultLocation,
			r.Store.LogRequest,
		},
		ResponseProcessors: []core.ResponseProcessor{
			processors.ConfirmMessage,
			r.Store.LogResponse,
		},
		Log:            r.Log,
		Telegram:       r.Services.Telegram,
		ClientRegistry: r.Store,
		SocketSender:   r.Services.SocketSender,
	})

	return r, nil
}

func (r *Runtime) Start() error {
	go r.Services.Push.Start(context.Background())
	go r.API.Start(context.Background())
	return http.ListenAndServe(r.Env.HttpHost, r.API)
}

func (r *Runtime) PushIntent(named string) core.PushIntentActor {
	for _, intent := range r.Intents {
		if pIntent, ok := intent.(core.PushIntentActor); ok && intent.IntentLabel() == named {
			return pIntent
		}
	}
	return nil
}
