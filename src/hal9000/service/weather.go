package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var ErrorWeatherForecastNotAvailable = errors.New("weather forecast not available")

type NOAAWeatherPointProperties struct {
	ForecastURL  string `json:"forecast"`
	RadarStation string `json:"radarStation"`
}

type NOAAWeatherPointResponse struct {
	Properties NOAAWeatherPointProperties `json:"properties"`
}

type NOAAWeatherForecastPeriod struct {
	StartTime        time.Time `json:"startTime"`
	EndTime          time.Time `json:"endTime"`
	DetailedForecast string    `json:"detailedForecast"`
}

type NOAAWeatherForecastProperties struct {
	Periods []NOAAWeatherForecastPeriod `json:"periods"`
}

type NOAAWeatherForecastResponse struct {
	Properties NOAAWeatherForecastProperties `json:"properties"`
}

func MakeWeatherAPICall(lat float64, lon float64, date time.Time) (string, string, error) {
	point, err := MakeWeatherAPIPointRequest(lat, lon)
	if err != nil {
		return "", "", nil
	}

	httpResponse, err := http.Get(point.ForecastURL)
	if err != nil {
		return "", "", nil
	}

	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return "", "", nil
	}

	var response NOAAWeatherForecastResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return "", "", nil
	}

	now := time.Now()
	for _, p := range response.Properties.Periods {
		if date.After(p.StartTime) && date.Before(p.EndTime) {
			radarURL := ""
			if now.After(p.StartTime) && now.Before(p.EndTime) {
				radarURL = fmt.Sprintf("https://radar.weather.gov/ridge/lite/%s_loop.gif", point.RadarStation)
			}
			return p.DetailedForecast, radarURL, nil
		}
	}

	return "", "", ErrorWeatherForecastNotAvailable
}

func MakeWeatherAPIPointRequest(lat float64, lon float64) (NOAAWeatherPointProperties, error) {
	httpResponse, err := http.Get(fmt.Sprintf("https://api.weather.gov/points/%f,%f", lat, lon))
	if err != nil {
		return NOAAWeatherPointProperties{}, err
	}

	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return NOAAWeatherPointProperties{}, err
	}

	var pointResponse NOAAWeatherPointResponse
	err = json.Unmarshal(responseBytes, &pointResponse)

	return pointResponse.Properties, err
}
