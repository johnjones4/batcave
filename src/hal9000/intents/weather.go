package intents

import (
	"fmt"
	"hal9000/types"
	"hal9000/util"
	"time"

	"github.com/codingsince1985/geo-golang/openstreetmap"
)

type weatherIntent struct {
	locale string
	date   time.Time
}

func NewWeatherIntent(m types.ParsedRequestMessage) (weatherIntent, error) {
	locale := ""
	for _, entity := range m.NamedEntities {
		if entity.Tag == util.NERTagPlace {
			locale = entity.Name
			break
		}
	}

	date := time.Now()
	if m.DateInfo != nil {
		date = m.DateInfo.Time
	}

	return weatherIntent{locale, date}, nil
}

func (i weatherIntent) Execute(runtime types.Runtime, lastState types.State) (types.State, types.ResponseMessage, error) {
	forecast, radarUrl, err := getWeather(runtime, i.date, i.locale)
	if err != nil {
		return lastState, types.ResponseMessage{}, err
	}

	responseMessage := formulateWeatherResponsePreamble(i.date, i.locale) + forecast
	return lastState, types.ResponseMessage{
		Text:  responseMessage,
		URL:   radarUrl,
		Extra: nil}, nil
}

func formulateWeatherResponsePreamble(date time.Time, locale string) string {
	message := "the weather"

	if locale != "" {
		message += fmt.Sprintf(" in %s", locale)
	}

	message += fmt.Sprintf(" for %s ", date.Format("January 2"))

	if date.After(time.Now().Add(time.Hour + 24)) {
		message += "will be: "
	} else {
		message += "is: "
	}

	return message
}

func getWeather(runtime types.Runtime, date time.Time, locale string) (string, string, error) {
	lat, lon := runtime.Weather().DefaultLatLon() //TODo location provider

	if locale != "" {
		geocoder := openstreetmap.Geocoder()
		location, err := geocoder.Geocode(locale)
		if err != nil {
			return "", "", err
		}
		lat = location.Lat
		lon = location.Lng
	}

	return runtime.Weather().MakeWeatherAPIForecastCall(lat, lon, date)
}
