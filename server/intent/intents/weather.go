package intents

import (
	"context"
	"fmt"
	"main/core"
	"main/services/noaa"
	"main/services/nominatim"
	"main/util"
	"time"

	"github.com/mitchellh/mapstructure"
)

type Weather struct {
	NOAA      *noaa.NOAA
	Nominatim *nominatim.Nominatim
}

type weatherIntentParseReceiver struct {
	Date     string `json:"date,omitempty"`
	Location string `json:"location,omitempty"`
}

func (p *Weather) IntentLabel() string {
	return "weather"
}

func (p *Weather) IntentParsePrompt(req *core.Request) string {
	return fmt.Sprintf("Extract the exact date and time relative to %s and location (use blank for unknown location) from the phrase \"%s\" and return the information in the JSON format {\"date\":\"RFC3339 format\", \"location\":\"city or town name, state, and country\"}", time.Now().String(), req.Message.Text)
}

func (p *Weather) IntentParseReceiver() any {
	return weatherIntentParseReceiver{}
}

func (p *Weather) ActOnIntent(ctx context.Context, req *core.Request, md *core.IntentMetadata) (core.Response, error) {
	var info weatherIntentParseReceiver
	err := mapstructure.Decode(md.IntentParseReceiver, &info)
	if err != nil {
		return core.ResponseEmpty, err
	}

	var location core.Coordinate
	var locationName string

	if info.Location != "" {
		location, err = p.Nominatim.Geocode(info.Location)
		if err != nominatim.ErrorLocationNotFound {
			if err != nil {
				return core.ResponseEmpty, err
			}
			locationName = info.Location
		}
	}

	if locationName == "" || location.Empty() {
		location = req.Coordinate
		locationName = "your location"
	}

	weather, err := p.NOAA.PredictWeather(location)
	if err != nil {
		return core.ResponseEmpty, err
	}

	idx := 0
	var media core.Media
	if info.Date != "" {
		parsedDate, err := util.ParseLLMDate(info.Date)
		if err != nil {
			return core.ResponseEmpty, err
		}
		for i, prediction := range weather.Forecast {
			if prediction.StartTime.Before(parsedDate) && prediction.EndTime.After(parsedDate) {
				idx = i
				break
			}
		}
	}
	if idx == 0 {
		media = core.Media{
			URL:  weather.RadarURL,
			Type: core.MediaTypeImage,
		}
	}

	return core.Response{
		OutboundMessage: core.OutboundMessage{
			Message: core.Message{
				Text: fmt.Sprintf("The weather from %s to %s in %s will be: %s",
					weather.Forecast[idx].StartTime.Format(core.FriendlyDateFormat),
					weather.Forecast[idx].EndTime.Format(core.FriendlyDateFormat),
					locationName,
					weather.Forecast[idx].DetailedForecast,
				),
			},
			Media: media,
		},
	}, nil
}
