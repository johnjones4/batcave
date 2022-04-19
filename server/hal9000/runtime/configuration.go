package runtime

import (
	"os"

	"github.com/johnjones4/hal-9000/server/hal9000/intent"
	"github.com/johnjones4/hal-9000/server/hal9000/learning"
	"github.com/johnjones4/hal-9000/server/hal9000/service"
	"github.com/johnjones4/hal-9000/server/hal9000/storage"
)

type Configuration struct {
	Abode          service.AbodeConfiguration
	Google         service.GoogleConfiguration
	Kasa           service.KasaConfiguration
	Metro          service.MetroConfiguration
	Trello         service.TrelloConfiguration
	WeatherStation service.WeatherStationConfiguration
	HouseProject   intent.HouseProjectConfiguration
	Storage        storage.Configuration
	Predictor      learning.PredictorConfiguration
}

func LoadConfigurationFromEnv() Configuration {
	return Configuration{
		Abode: service.AbodeConfiguration{
			Username: os.Getenv("ABODE_USERNAME"),
			Password: os.Getenv("ABODE_PASSWORD"),
		},
		Google: service.GoogleConfiguration{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RefreshToken: os.Getenv("GOOGLE_REFRESH_TOKEN"),
		},
		Kasa: service.KasaConfiguration{
			DevicesPath: os.Getenv("KASA_DEVICES_FILE"),
			MQTTURL:     os.Getenv("KASA_MQTT_URL"),
		},
		Metro: service.MetroConfiguration{
			APIKey: os.Getenv("METRO_API_KEY"),
		},
		Trello: service.TrelloConfiguration{
			APIKey: os.Getenv("TRELLO_API_KEY"),
			Token:  os.Getenv("TRELLO_TOKEN"),
		},
		WeatherStation: service.WeatherStationConfiguration{
			Upstream: os.Getenv("WEATHER_STATION_UPSTREAM"),
		},
		HouseProject: intent.HouseProjectConfiguration{
			ListId: os.Getenv("TRELLO_LID_HOUSE_TODO"),
		},
		Storage: storage.Configuration{
			ClientsPath: os.Getenv("CLIENT_STORE_FILE"),
			UsersPath:   os.Getenv("USER_STORE_FILE"),
			DatabaseURL: os.Getenv("DATABASE_URL"),
		},
		Predictor: learning.PredictorConfiguration{
			IntentMapFilePath: os.Getenv("INTENT_MAP_FILE"),
			ModelFilePath:     os.Getenv("MODEL_FILE"),
		},
	}
}
