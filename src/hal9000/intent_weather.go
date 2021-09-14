package hal9000

import (
	"fmt"
	"hal9000/service"
	"hal9000/util"
	"os"
	"strconv"
	"time"

	"github.com/codingsince1985/geo-golang/openstreetmap"
)

type WeatherIntent struct {
	Locale string
	Date   time.Time
}

func NewWeatherIntent(m ParsedRequestMessage) (WeatherIntent, error) {
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

	return WeatherIntent{locale, date}, nil
}

func (i WeatherIntent) Execute(lastState State) (State, util.ResponseMessage, error) {
	forecast, radarUrl, err := GetWeather(i.Date, i.Locale)
	if err != nil {
		return lastState, util.ResponseMessage{}, err
	}

	responseMessage := FormulateWeatherResponsePreamble(i.Date, i.Locale) + forecast
	return lastState, util.ResponseMessage{
		Text:  responseMessage,
		URL:   radarUrl,
		Extra: nil}, nil
}

func FormulateWeatherResponsePreamble(date time.Time, locale string) string {
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

func GetWeather(date time.Time, locale string) (string, string, error) {
	lat, err := strconv.ParseFloat(os.Getenv("DEFAULT_WEATHER_LAT"), 64)
	if err != nil {
		return "", "", err
	}
	lon, err := strconv.ParseFloat(os.Getenv("DEFAULT_WEATHER_LON"), 64)
	if err != nil {
		return "", "", err
	}

	if locale != "" {
		geocoder := openstreetmap.Geocoder()
		location, err := geocoder.Geocode(locale)
		if err != nil {
			return "", "", err
		}
		lat = location.Lat
		lon = location.Lng
	}

	return service.MakeWeatherAPIForecastCall(lat, lon, date)
}
