package api

import (
	"net/http"

	"github.com/johnjones4/hal-9000/server/hal9000/core"
	"github.com/johnjones4/hal-9000/server/hal9000/intent"
	"github.com/johnjones4/hal-9000/server/hal9000/learning"
	"github.com/johnjones4/hal-9000/server/hal9000/security"
	"github.com/johnjones4/hal-9000/server/hal9000/service"
	"github.com/johnjones4/hal-9000/server/hal9000/storage"
)

func New(userStoreFile, stateStoreFile, logFile, tokenKey string) (http.Handler, error) {
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
	h := intent.IntentSet{
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

	us := storage.NewUserStore(userStoreFile)
	err = us.Load()
	if err != nil {
		return nil, err
	}

	ss := storage.NewStateStore(stateStoreFile)
	err = ss.Load()
	if err != nil {
		return nil, err
	}

	logger, err := learning.NewInteractionLogger(logFile)
	if err != nil {
		return nil, err
	}

	tm, err := security.NewTokenManager([]byte(tokenKey))
	if err != nil {
		return nil, err
	}

	router := makeServer(&h, us, ss, logger, tm)
	return router, nil

}
