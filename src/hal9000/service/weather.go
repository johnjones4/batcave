package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"hal9000/util"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var ErrorWeatherForecastNotAvailable = errors.New("weather forecast not available")

type NOAAWeatherPointProperties struct {
	ForecastURL  string `json:"forecast"`
	ForecastZone string `json:"forecastZone"`
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

type NOAAWeatherAlertFeatureProperties struct {
	ID            string   `json:"id"`
	AffectedZones []string `json:"affectedZones"`
	Headline      string   `json:"headline"`
}

type NOAAWeatherAlertFeature struct {
	Properties NOAAWeatherAlertFeatureProperties `json:"properties"`
}

type NOAAWeatherAlertResponse struct {
	Features []NOAAWeatherAlertFeature `json:"features"`
}

func StartWeatherAlertLoop(alertChan *chan util.ResponseMessage) {
	lat, err := strconv.ParseFloat(os.Getenv("DEFAULT_WEATHER_LAT"), 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	lon, err := strconv.ParseFloat(os.Getenv("DEFAULT_WEATHER_LON"), 64)
	if err != nil {
		fmt.Println(err)
		return
	}

	previouslyHandledAlerts := util.NewStringBuffer(1000)

	for {
		_, radar, _ := MakeWeatherAPIForecastCall(lat, lon, time.Now())
		alerts, err := MakeWeatherAPIAlertCall(lat, lon)
		if err != nil {
			fmt.Println(err)
		} else {
			for _, alert := range alerts {
				if !previouslyHandledAlerts.Contains(alert.Properties.ID) {
					*alertChan <- util.ResponseMessage{Text: fmt.Sprintf("Weather alert: %s", alert.Properties.Headline), URL: radar, Extra: nil}
					previouslyHandledAlerts.Push(alert.Properties.ID)
				}
			}
		}
		time.Sleep(time.Hour / 2)
	}
}

func MakeWeatherAPIAlertCall(lat float64, lon float64) ([]NOAAWeatherAlertFeature, error) {
	point, err := MakeWeatherAPIPointRequest(lat, lon)
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

	var response NOAAWeatherAlertResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return nil, err
	}

	features := make([]NOAAWeatherAlertFeature, 0)
	for _, feature := range response.Features {
		if util.ContainsString(feature.Properties.AffectedZones, point.ForecastZone) {
			features = append(features, feature)
		}

	}

	return features, nil
}

func MakeWeatherAPIForecastCall(lat float64, lon float64, date time.Time) (string, string, error) {
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

	for i, p := range response.Properties.Periods {
		if date.After(p.StartTime) && date.Before(p.EndTime) {
			radarURL := ""
			if i == 0 {
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
