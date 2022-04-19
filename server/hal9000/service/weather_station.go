package service

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type WeatherStationConfiguration struct {
	Upstream string
}

type WeatherStation struct {
	configuration WeatherStationConfiguration
}

type WeatherStatonResponse struct {
	Timestamp        time.Time `json:"timestamp"`
	Uptime           int64     `json:"uptime"`
	AvgWindSpeed     float64   `json:"avg_wind_speed"`
	MinWindSpeed     float64   `json:"min_wind_speed"`
	MaxWindSpeed     float64   `json:"max_wind_speed"`
	Temperature      float64   `json:"temperature"`
	Gas              float64   `json:"gas"`
	RelativeHumidity float64   `json:"relative_humidity"`
	Pressure         float64   `json:"pressure"`
}

func NewWeatherStation(configuration WeatherStationConfiguration) *WeatherStation {
	return &WeatherStation{
		configuration: configuration,
	}
}

func (w *WeatherStation) GetWeather() (WeatherStatonResponse, error) {
	res, err := http.Get(w.configuration.Upstream)
	if err != nil {
		return WeatherStatonResponse{}, nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return WeatherStatonResponse{}, nil
	}

	var info WeatherStatonResponse
	err = json.Unmarshal(body, &info)
	if err != nil {
		return WeatherStatonResponse{}, nil
	}

	return info, nil
}
