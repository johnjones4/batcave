package runtime

import (
	"os"

	"github.com/johnjones4/hal-9000/server/hal9000/core"
	"github.com/johnjones4/hal-9000/server/hal9000/intent"
	"github.com/johnjones4/hal-9000/server/hal9000/learning"
	"github.com/johnjones4/hal-9000/server/hal9000/security"
	"github.com/johnjones4/hal-9000/server/hal9000/service"
	"github.com/johnjones4/hal-9000/server/hal9000/storage"
)

type Runtime struct {
	Intents      *IntentSet
	UserStore    *storage.UserStore
	StateStore   *storage.StateStore
	Logger       *learning.InteractionLogger
	TokenManager *security.TokenManager
}

func New() (*Runtime, error) {
	kasa, err := service.NewKasa()
	if err != nil {
		return nil, err
	}
	intents := []core.Intent{
		&intent.Forecast{
			Service: service.NewNOAA(),
		},
		&intent.Metro{
			Service: service.NewMetro(),
		},
		&intent.Schedule{
			Service: service.NewGoogle(),
		},
		&intent.WeatherStation{
			Service: service.NewWeatherStation(),
		},
		&intent.Lights{
			Service: kasa,
		},
		&intent.Abode{
			Service: service.NewAbode(),
		},
	}
	h := &IntentSet{
		Intents: append([]core.Intent{
			&intent.Forecast{
				Service: service.NewNOAA(),
			},
			&intent.Metro{
				Service: service.NewMetro(),
			},
			&intent.Schedule{
				Service: service.NewGoogle(),
			},
			&intent.WeatherStation{
				Service: service.NewWeatherStation(),
			},
			&intent.Lights{
				Service: kasa,
			},
			&intent.Abode{
				Service: service.NewAbode(),
			},
		}, &intent.Info{
			Intents: intents,
		}),
	}

	us := storage.NewUserStore(os.Getenv("USER_STORE_FILE"))
	err = us.Load()
	if err != nil {
		return nil, err
	}

	ss := storage.NewStateStore(os.Getenv("STATE_STORE_FILE"))
	err = ss.Load()
	if err != nil {
		return nil, err
	}

	logger, err := learning.NewInteractionLogger(os.Getenv("LOG_FILE"))
	if err != nil {
		return nil, err
	}

	tm, err := security.NewTokenManager([]byte(os.Getenv("TOKEN_KEY")))
	if err != nil {
		return nil, err
	}

	return &Runtime{
		Intents:      h,
		UserStore:    us,
		StateStore:   ss,
		Logger:       logger,
		TokenManager: tm,
	}, nil
}
