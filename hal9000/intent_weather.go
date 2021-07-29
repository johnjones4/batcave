package hal9000

import (
	"errors"
	"fmt"
	"hal9000/service"
	"hal9000/util"
	"os"
	"strconv"
	"time"

	"github.com/codingsince1985/geo-golang/openstreetmap"
)

type WeatherIntent struct {
	Locale string    `json:"locale"`
	Date   time.Time `json:"date"`
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

func (i WeatherIntent) Execute(lastState State) (State, ResponseMessage, error) {
	forecast, err := GetWeather(i.Date, i.Locale)
	if err != nil {
		return lastState, ResponseMessage{}, err
	}

	responseMessage := FormulateWeatherResponsePreamble(i.Date, i.Locale) + forecast
	return lastState, ResponseMessage{responseMessage, "", nil}, nil
}

func FormulateWeatherResponsePreamble(date time.Time, locale string) string {
	message := "the weather"

	if locale != "" {
		message += fmt.Sprintf(" in %s", locale)
	}

	message += fmt.Sprintf(" for %s ", date.Format("July _2"))

	if date.After(time.Now().Add(time.Hour + 24)) {
		message += "will be: "
	} else {
		message += "is: "
	}

	return message
}

func GetWeather(date time.Time, locale string) (string, error) {
	lat, err := strconv.ParseFloat(os.Getenv("DEFAULT_WEATHER_LAT"), 64)
	if err != nil {
		return "", err
	}
	lon, err := strconv.ParseFloat(os.Getenv("DEFAULT_WEATHER_LON"), 64)
	if err != nil {
		return "", err
	}

	if locale != "" {
		geocoder := openstreetmap.Geocoder()
		location, err := geocoder.Geocode(locale)
		if err != nil {
			return "", err
		}
		lat = location.Lat
		lon = location.Lng
	}

	response, err := service.MakeWeatherAPICall(lat, lon)
	if err != nil {
		return "", err
	}

	if date.Before(time.Now().Add(time.Hour + 24)) {
		respString := WeatherDetailsToString(response.Current.Weather)
		respString += fmt.Sprintf("The temperature is %.0f degrees, feeling like %.0f. ", response.Current.Temperature, response.Current.FeelsLike)
		respString += fmt.Sprintf("The humidity is %.0f%% and the dew point is at %.0f.", response.Current.Humidity, response.Current.DewPoint)
		return respString, nil
	}

	for _, day := range response.Daily {
		if int(date.Unix()) > day.Timestamp {
			respString := WeatherDetailsToString(day.Weather)
			respString += fmt.Sprintf("The high will be %.0f degrees, the low will be %0.f degrees, and by mid-day it will feel like %.0f. ", day.Temperature.Max, day.Temperature.Min, day.FeelsLike.Day)
			respString += fmt.Sprintf("The humidity will be %.0f%% and the dew point will be at %.0f.", day.Humidity, day.DewPoint)
			return respString, nil
		}
	}

	return "", errors.New("could not parse weather")
}

func WeatherDetailsToString(weather []service.WeatherAPIResponseWeatherDetail) string {
	if len(weather) == 1 {
		return weather[0].Description + ". "
	} else if len(weather) > 1 {
		respString := ""
		for i, item := range weather {
			if i == len(weather)-1 {
				respString += ", and "
			} else if i > 0 {
				respString += ", "
			}
			respString += item.Description
		}
		respString += ". "
		return respString
	} else {
		return ""
	}
}
