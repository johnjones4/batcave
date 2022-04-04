package service

import (
	"encoding/json"
	"io"
	"main/core"
	"math"
	"net/http"
	"os"
)

const (
	urlRoot = "https://api.wmata.com"
)

type Metro struct {
	apiKey string
}

func NewMetro() *Metro {
	return &Metro{
		apiKey: os.Getenv("METRO_API_KEY"),
	}
}

func (m Metro) GetArrivals(c core.Coordinate) ([]Arrival, error) {
	station, err := m.findClosestStation(c)
	if err != nil {
		return nil, err
	}

	return m.getArrivalsForStation(station)
}

type stationResponse struct {
	Stations []struct {
		Lat  float64 `json:"Lat"`
		Lon  float64 `json:"Lon"`
		Code string  `json:"Code"`
	} `json:"Stations"`
}

func (m Metro) findClosestStation(c core.Coordinate) (string, error) {
	req, err := http.NewRequest("GET", urlRoot+"/Rail.svc/json/jStations", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("api_key", m.apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var info stationResponse
	err = json.Unmarshal(body, &info)
	if err != nil {
		return "", err
	}

	closestDistance := math.MaxFloat64
	closesStation := ""
	for _, station := range info.Stations {
		d := c.DistanceTo(core.Coordinate{
			Latitude:  station.Lat,
			Longitude: station.Lon,
		})
		if d < float64(closestDistance) {
			closestDistance = d
			closesStation = station.Code
		}
	}
	return closesStation, nil
}

type arrivalsResponse struct {
	Trains []Arrival `json:"Trains"`
}

type Arrival struct {
	Destination string `json:"DestinationName"`
	Line        string `json:"Line"`
	Min         string `json:"Min"`
}

func (m Metro) getArrivalsForStation(code string) ([]Arrival, error) {
	req, err := http.NewRequest("GET", urlRoot+"/StationPrediction.svc/json/GetPrediction/"+code, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("api_key", m.apiKey)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var info arrivalsResponse
	err = json.Unmarshal(body, &info)
	if err != nil {
		return nil, err
	}

	return info.Trains, nil
}
