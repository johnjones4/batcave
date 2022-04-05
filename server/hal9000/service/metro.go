package service

import (
	"encoding/json"
	"io"
	"math"
	"net/http"
	"os"

	"github.com/johnjones4/hal-9000/server/hal9000/core"
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

func (m Metro) GetArrivals(c core.Coordinate) (Station, []Arrival, error) {
	station, err := m.findClosestStation(c)
	if err != nil {
		return Station{}, nil, err
	}

	arrivals, err := m.getArrivalsForStation(station)
	if err != nil {
		return Station{}, nil, err
	}

	return station, arrivals, nil
}

type Station struct {
	Lat  float64 `json:"Lat"`
	Lon  float64 `json:"Lon"`
	Code string  `json:"Code"`
	Name string  `json:"Name"`
}

type stationResponse struct {
	Stations []Station `json:"Stations"`
}

func (m Metro) findClosestStation(c core.Coordinate) (Station, error) {
	req, err := http.NewRequest("GET", urlRoot+"/Rail.svc/json/jStations", nil)
	if err != nil {
		return Station{}, err
	}

	req.Header.Set("api_key", m.apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Station{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Station{}, err
	}

	var info stationResponse
	err = json.Unmarshal(body, &info)
	if err != nil {
		return Station{}, err
	}

	closestDistance := math.MaxFloat64
	var closesStation Station
	for _, station := range info.Stations {
		d := c.DistanceTo(core.Coordinate{
			Latitude:  station.Lat,
			Longitude: station.Lon,
		})
		if d < float64(closestDistance) {
			closestDistance = d
			closesStation = station
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

func (m Metro) getArrivalsForStation(station Station) ([]Arrival, error) {
	req, err := http.NewRequest("GET", urlRoot+"/StationPrediction.svc/json/GetPrediction/"+station.Code, nil)
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
