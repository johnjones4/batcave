package noaa

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"main/core"
	"net/http"
	"strings"
	"time"
)

type noaaWeatherPointProperties struct {
	ForecastURL  string `json:"forecast"`
	ForecastZone string `json:"forecastZone"`
	RadarStation string `json:"radarStation"`
}

type noaaWeatherPointResponse struct {
	Properties noaaWeatherPointProperties `json:"properties"`
}

type noaaWeatherForecastProperties struct {
	Periods []core.WeatherForecastPeriod `json:"periods"`
}

type noaaWeatherForecastResponse struct {
	Properties noaaWeatherForecastProperties `json:"properties"`
}

type noaaWeatherAlertFeature struct {
	Properties core.WeatherAlert `json:"properties"`
}

type noaaWeatherAlertResponse struct {
	Features []noaaWeatherAlertFeature `json:"features"`
}

type NOAA struct {
}

func (f *NOAA) PredictWeather(coord core.Coordinate) (core.WeatherForecast, error) {
	point, err := makeWeatherAPIPointRequest(coord)
	if err != nil {
		return core.WeatherForecast{}, err
	}

	radarURL := fmt.Sprintf("https://radar.weather.gov/ridge/standard/%s_loop.gif?refreshed=%d", point.RadarStation, time.Now().Unix())

	forecast, err := makeWeatherAPIForecastCall(point)
	if err != nil {
		return core.WeatherForecast{}, err
	}

	alerts, err := makeWeatherAPIAlertCall(point)
	if err != nil {
		return core.WeatherForecast{}, err
	}

	return core.WeatherForecast{
		RadarURL: radarURL,
		Forecast: forecast,
		Alerts:   alerts,
	}, nil
}

func makeWeatherAPIPointRequest(coord core.Coordinate) (noaaWeatherPointProperties, error) {
	url := fmt.Sprintf("https://api.weather.gov/points/%f,%f", coord.Latitude, coord.Longitude)
	httpResponse, err := http.Get(url)
	if err != nil {
		return noaaWeatherPointProperties{}, err
	}

	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return noaaWeatherPointProperties{}, err
	}

	var pointResponse noaaWeatherPointResponse
	err = json.Unmarshal(responseBytes, &pointResponse)

	return pointResponse.Properties, err
}

func makeWeatherAPIForecastCall(point noaaWeatherPointProperties) ([]core.WeatherForecastPeriod, error) {
	httpResponse, err := http.Get(point.ForecastURL)
	if err != nil {
		return nil, err
	}

	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	var response noaaWeatherForecastResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Properties.Periods) == 0 {
		return nil, errors.New("no forecast data returned")
	}

	return response.Properties.Periods, nil
}

func makeWeatherAPIAlertCall(point noaaWeatherPointProperties) ([]core.WeatherAlert, error) {
	zoneId := strings.Replace(point.ForecastZone, "https://api.weather.gov/zones/forecast/", "", 1)

	httpResponse, err := http.Get(fmt.Sprintf("https://api.weather.gov/alerts/active/zone/%s", zoneId))
	if err != nil {
		return nil, err
	}

	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	var response noaaWeatherAlertResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return nil, err
	}

	featureProps := make([]core.WeatherAlert, 0)
	for _, feature := range response.Features {
		for _, zone := range feature.Properties.AffectedZones {
			if zone == point.ForecastZone {
				featureProps = append(featureProps, feature.Properties)
			}
		}
	}
	return featureProps, nil
}
