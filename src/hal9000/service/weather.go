package service

import (
	"encoding/json"
	"fmt"
	"hal9000/types"
	"hal9000/util"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type weatherProviderConcrete struct {
	defaultLat float64
	defaultLon float64
}

func InitWeatherProvider(runtime *types.Runtime) (types.WeatherProvider, error) {
	lat, err := strconv.ParseFloat(os.Getenv("DEFAULT_WEATHER_LAT"), 64)
	if err != nil {
		return nil, err
	}
	lon, err := strconv.ParseFloat(os.Getenv("DEFAULT_WEATHER_LON"), 64)
	if err != nil {
		return nil, err
	}
	wp := weatherProviderConcrete{lat, lon}
	go (func() {
		previouslyHandledAlerts := util.NewStringBuffer(1000)
		for {
			_, radar, _ := wp.MakeWeatherAPIForecastCall(lat, lon, time.Now())
			alerts, err := wp.MakeWeatherAPIAlertCall(lat, lon)
			if err != nil {
				(*(*runtime).Logger()).LogError(err)
			} else {
				for _, alert := range alerts {
					if !previouslyHandledAlerts.Contains(alert.Properties.ID) {
						m := types.ResponseMessage{Text: fmt.Sprintf("Weather alert: %s", alert.Properties.Headline), URL: radar, Extra: nil}
						(*(*runtime).AlertQueue()).Enqueue(m)
						previouslyHandledAlerts.Push(alert.Properties.ID)
					}
				}
			}
			time.Sleep(time.Hour / 2)
		}
	})()
	return wp, nil
}

func (wp weatherProviderConcrete) DefaultLatLon() (float64, float64) {
	return wp.defaultLat, wp.defaultLon
}

func (wp weatherProviderConcrete) MakeWeatherAPIAlertCall(lat float64, lon float64) ([]types.NOAAWeatherAlertFeature, error) {
	point, err := wp.MakeWeatherAPIPointRequest(lat, lon)
	if err != nil {
		return nil, err
	}

	httpResponse, err := http.Get(fmt.Sprintf("https://api.weather.gov/alerts/active/area/%s", os.Getenv("DEFAULT_WEATHER_ALERT_AREA")))
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

	features := make([]types.NOAAWeatherAlertFeature, 0)
	for _, feature := range response.Features {
		if util.ContainsString(feature.Properties.AffectedZones, point.ForecastZone) {
			features = append(features, feature)
		}

	}

	return features, nil
}

func (wp weatherProviderConcrete) MakeWeatherAPIForecastCall(lat float64, lon float64, date time.Time) (string, string, error) {
	point, err := wp.MakeWeatherAPIPointRequest(lat, lon)
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

	var response types.NOAAWeatherForecastResponse
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

	return "", "", util.ErrorWeatherForecastNotAvailable
}

func (wp weatherProviderConcrete) MakeWeatherAPIPointRequest(lat float64, lon float64) (types.NOAAWeatherPointProperties, error) {
	httpResponse, err := http.Get(fmt.Sprintf("https://api.weather.gov/points/%f,%f", lat, lon))
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
