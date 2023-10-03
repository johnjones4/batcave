package intent

import (
	"main/core"
	"main/intent/intents"
	"main/services"
)

func Factory(services *services.Services) []core.IntentActor {
	return []core.IntentActor{
		&intents.Play{
			TuneIn: services.TuneIn,
			Push:   services.Push,
		},
		&intents.Stop{},
		&intents.ToggleDevice{
			HomeAssistant: services.HomeAssistant,
		},
		&intents.Weather{
			Weather:  services.NOAA,
			Geocoder: services.Nominatim,
		},
		&intents.Remind{
			Push: services.Push,
		},
		&intents.Unknown{},
	}
}
