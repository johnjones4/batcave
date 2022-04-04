package main

import (
	"main/core"
	"main/intent"
	"main/learning"
	"main/security"
	"main/service"
	"main/storage"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	kasa, err := service.NewKasa()
	if err != nil {
		panic(err)
	}
	h := intent.IntentSet{
		Intents: []core.Intent{
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
		},
	}

	us := storage.NewUserStore(os.Getenv("USER_STORE_FILE"))
	err = us.Load()
	if err != nil {
		panic(err)
	}

	ss := storage.NewStateStore(os.Getenv("STATE_STORE_FILE"))
	err = ss.Load()
	if err != nil {
		panic(err)
	}

	logger, err := learning.NewInteractionLogger(os.Getenv("LOG_FILE"))
	if err != nil {
		panic(err)
	}

	tm, err := security.NewTokenManager([]byte(os.Getenv("TOKEN_KEY")))
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(os.Getenv("HTTP_HOST"), makeServer(&h, us, ss, logger, tm))
	panic(err)
}
