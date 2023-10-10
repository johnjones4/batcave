package nominatim

import (
	"encoding/json"
	"errors"
	"io"
	"main/core"
	"net/http"
	"net/url"
	"strconv"
)

var (
	ErrorLocationNotFound = errors.New("could not determine location")
)

type Nominatim struct {
}

type nominatimResponseItem struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"lon"`
	Address   struct {
		CountryCode string `json:"country_code"`
	} `json:"address"`
}

func (n *Nominatim) Geocode(q string) (core.Coordinate, error) {
	params := url.Values{
		"q":              []string{q},
		"format":         []string{"json"},
		"addressdetails": []string{"1"},
		"limit":          []string{"1"},
	}

	url := "https://nominatim.openstreetmap.org/search?" + params.Encode()

	res, err := http.Get(url)
	if err != nil {
		return core.Coordinate{}, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return core.Coordinate{}, err
	}

	var items []nominatimResponseItem
	err = json.Unmarshal(resBody, &items)
	if err != nil {
		return core.Coordinate{}, err
	}

	if len(items) == 0 || items[0].Address.CountryCode != "us" {
		return core.Coordinate{}, ErrorLocationNotFound
	}

	lat, err := strconv.ParseFloat(items[0].Latitude, 64)
	if err != nil {
		return core.Coordinate{}, err
	}

	lon, err := strconv.ParseFloat(items[0].Longitude, 64)
	if err != nil {
		return core.Coordinate{}, err
	}

	return core.Coordinate{
		Latitude:  lat,
		Longitude: lon,
	}, nil
}
