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
	kasa, err := service.NewKasa(os.Getenv("KASA_DEVICES_FILE"), os.Getenv("KASA_MQTT_URL"))
	if err != nil {
		return nil, err
	}
	intents := []core.Intent{
		&intent.Forecast{
			Service: service.NewNOAA(),
		},
		&intent.Metro{
			Service: service.NewMetro(os.Getenv("METRO_API_KEY")),
		},
		&intent.Schedule{
			Service: service.NewGoogle(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), os.Getenv("GOOGLE_REFRESH_TOKEN")),
		},
		&intent.WeatherStation{
			Service: service.NewWeatherStation(os.Getenv("WEATHER_STATION_UPSTREAM")),
		},
		&intent.Lights{
			Service: kasa,
		},
		&intent.Abode{
			Service: service.NewAbode(os.Getenv("ABODE_USERNAME"), os.Getenv("ABODE_PASSWORD")),
		},
		&intent.HouseProject{
			Service: service.NewTrello(os.Getenv("TRELLO_API_KEY"), os.Getenv("TRELLO_TOKEN")),
			ListId:  os.Getenv("TRELLO_LID_HOUSE_TODO"),
		},
	}
	h := &IntentSet{
		Intents: append(intents, &intent.Info{
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
