package runtime

import (
	"github.com/johnjones4/hal-9000/server/hal9000/core"
	"github.com/johnjones4/hal-9000/server/hal9000/intent"
	"github.com/johnjones4/hal-9000/server/hal9000/learning"
	"github.com/johnjones4/hal-9000/server/hal9000/service"
	"github.com/johnjones4/hal-9000/server/hal9000/storage"
)

type Runtime struct {
	Intents     *IntentSet
	StateStore  *storage.StateStore
	Logger      *storage.InteractionLogger
	ClientStore *storage.ClientStore
	UserStore   *storage.UserStore
	Predictor   *learning.Predictor
}

func New() (*Runtime, error) {
	configuration := LoadConfigurationFromEnv()

	kasa, err := service.NewKasa(configuration.Kasa)
	if err != nil {
		return nil, err
	}
	intents := []core.Intent{
		&intent.Forecast{
			Service: service.NewNOAA(),
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
	}
	h := &IntentSet{
		Intents: append(intents, &intent.Info{
			Intents: intents,
		}),
	}

	pool, err := storage.Connect(configuration.Storage)
	if err != nil {
		return nil, err
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

	predictor, err := learning.NewPredictor(configuration.Predictor)
	if err != nil {
		return nil, err
	}

	return &Runtime{
		Intents:     h,
		StateStore:  storage.NewStateStore(pool),
		Logger:      storage.NewInteractionLogger(pool),
		ClientStore: cs,
		UserStore:   us,
		Predictor:   predictor,
	}, nil
}
