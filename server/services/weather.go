package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/types"
	"net/http"
	"strings"
	"time"
)

func StartWeatherUpdater(coords []types.Coordinate, updateChan chan []types.Weather, errorChan chan error) {
	for {
		weathers := make([]types.Weather, 0)
		for _, coord := range coords {
			weather, err := getWeather(coord)
			if err != nil {
				errorChan <- err
			} else {
				weathers = append(weathers, weather)
			}
		}
		updateChan <- weathers
		time.Sleep(time.Minute * 5)
	}
}

func getWeather(coord types.Coordinate) (types.Weather, error) {
	point, err := makeWeatherAPIPointRequest(coord)
	if err != nil {
		return types.Weather{}, nil
	}

	radarURL := fmt.Sprintf("https://radar.weather.gov/ridge/lite/%s_loop.gif?v=%d", point.RadarStation, time.Now().Unix())

	forecast, err := makeWeatherAPIForecastCall(point)
	if err != nil {
		return types.Weather{}, nil
	}

	alerts, err := makeWeatherAPIAlertCall(point)
	if err != nil {
		return types.Weather{}, nil
	}

	return types.Weather{
		RadarURL: radarURL,
		Forecast: forecast,
		Alerts:   alerts,
	}, nil
}

func makeWeatherAPIPointRequest(coord types.Coordinate) (types.NOAAWeatherPointProperties, error) {
	httpResponse, err := http.Get(fmt.Sprintf("https://api.weather.gov/points/%f,%f", coord.Latitude, coord.Longitude))
	if err != nil {
		return types.NOAAWeatherPointProperties{}, err
	}

	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return types.NOAAWeatherPointProperties{}, err
	}

	var pointResponse types.NOAAWeatherPointResponse
	err = json.Unmarshal(responseBytes, &pointResponse)

	return pointResponse.Properties, err
}

func makeWeatherAPIForecastCall(point types.NOAAWeatherPointProperties) ([]types.NOAAWeatherForecastPeriod, error) {
	httpResponse, err := http.Get(point.ForecastURL)
	if err != nil {
		return nil, nil
	}

	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, nil
	}

	var response types.NOAAWeatherForecastResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return nil, nil
	}

	return response.Properties.Periods, nil
}

func makeWeatherAPIAlertCall(point types.NOAAWeatherPointProperties) ([]types.NOAAWeatherAlertFeatureProperties, error) {
	zoneId := strings.Replace(point.ForecastZone, "https://api.weather.gov/zones/forecast/", "", 1)

	httpResponse, err := http.Get(fmt.Sprintf("https://api.weather.gov/alerts/active/zone/%s", zoneId))
	if err != nil {
		return nil, err
	}

	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	var response types.NOAAWeatherAlertResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return nil, err
	}

	featureProps := make([]types.NOAAWeatherAlertFeatureProperties, 0)
	for _, feature := range response.Features {
		for _, zone := range feature.Properties.AffectedZones {
			if zone == point.ForecastZone {
				featureProps = append(featureProps, feature.Properties)
			}
		}
	}
	return featureProps, nil
}
