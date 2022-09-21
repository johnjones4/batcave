package runtime

import (
	"github.com/johnjones4/hal-9000/server/hal9000/core"
	"github.com/johnjones4/hal-9000/server/hal9000/intent"
	"github.com/johnjones4/hal-9000/server/hal9000/learning"
	"github.com/johnjones4/hal-9000/server/hal9000/service"
	"github.com/johnjones4/hal-9000/server/hal9000/storage"
)

type Runtime struct {
	Intents          *IntentSet
	StateStore       storage.StateStore
	Logger           storage.InteractionLogger
	ClientStore      *storage.ClientStore
	UserStore        *storage.UserStore
	IntentPredictor  *learning.IntentPredictor
	VoiceTranscriber *learning.VoiceTranscriber
}

func New() (*Runtime, error) {
	configuration := LoadConfigurationFromEnv()
	var err error

	kasa, err := service.NewKasa(configuration.Kasa)
	if err != nil {
		return nil, err
	}
	nest, err := service.NewNestCameras(configuration.NestCameras)
	if err != nil {
		return nil, err
	}
	noaa := service.NewNOAA()
	intents := []core.Intent{
		&intent.Forecast{
			Service: noaa,
		},
		&intent.Metro{
			Service: service.NewMetro(configuration.Metro),
		},
		&intent.Schedule{
			Service: service.NewGoogle(configuration.Google),
		},
		&intent.WeatherStation{
			Service: service.NewWeatherStation(configuration.WeatherStation),
		},
		&intent.Lights{
			Service: kasa,
		},
		&intent.Abode{
			Service: service.NewAbode(configuration.Abode),
		},
		&intent.HouseProject{
			Service:       service.NewTrello(configuration.Trello),
			Configuration: configuration.HouseProject,
		},
		&intent.Display{
			DisplayServices: []service.DisplayService{
				nest,
				noaa,
			},
		},
	}
	h := &IntentSet{
		Intents: append(intents, &intent.Info{
			Intents: intents,
		}),
	}

	var logger storage.InteractionLogger
	var stateStore storage.StateStore
	if configuration.Storage.DatabaseURL != "" {
		pool, err := storage.Connect(configuration.Storage)
		if err != nil {
			return nil, err
		}
		logger = storage.NewDatabaseInteractionLogger(pool)
		stateStore = storage.NewDatabaseStateStore(pool)
	} else {
		logger = &storage.TerminalInteractionLogger{}
		stateStore = storage.NewMemoryStateStore()
	}

	cs := storage.NewClientStore(configuration.Storage)
	err = cs.Load()
	if err != nil {
		return nil, err
	}

	us := storage.NewUserStore(configuration.Storage)
	err = us.Load()
	if err != nil {
		return nil, err
	}

	predictor, err := learning.NewIntentPredictor(configuration.IntentPredictor)
	if err != nil {
		return nil, err
	}

	transcriber, err := learning.NewVoiceTranscriber(configuration.VoiceTranscriber)
	if err != nil {
		return nil, err
	}

	return &Runtime{
		Intents:          h,
		StateStore:       stateStore,
		Logger:           logger,
		ClientStore:      cs,
		UserStore:        us,
		IntentPredictor:  predictor,
		VoiceTranscriber: transcriber,
	}, nil
}
