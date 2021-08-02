package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type WeatherAPIResponseWeatherDetail struct {
	Description string `json:"description"`
}

type WeatherAPIResponseCurrent struct {
	Temperature float64                           `json:"temp"`
	FeelsLike   float64                           `json:"feels_like"`
	Humidity    float64                           `json:"humidity"`
	DewPoint    float64                           `json:"dew_point"`
	Weather     []WeatherAPIResponseWeatherDetail `json:"weather"`
}

type WeatherAPIResponseForecastTemp struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type WeatherAPIResponseForecastFeelsLike struct {
	Day     float64 `json:"day"`
	Night   float64 `json:"night"`
	Evening float64 `json:"eve"`
	Morning float64 `json:"morn"`
}

type WeatherAPIResponseForecast struct {
	Temperature WeatherAPIResponseForecastTemp      `json:"temp"`
	FeelsLike   WeatherAPIResponseForecastFeelsLike `json:"feels_like"`
	Humidity    float64                             `json:"humidity"`
	DewPoint    float64                             `json:"dew_point"`
	Weather     []WeatherAPIResponseWeatherDetail   `json:"weather"`
	Timestamp   int                                 `json:"dt"`
}

type WeatherAPIResponse struct {
	Current WeatherAPIResponseCurrent    `json:"current"`
	Daily   []WeatherAPIResponseForecast `json:"daily"`
}

func MakeWeatherAPICall(lat float64, lon float64) (WeatherAPIResponse, error) {
	apiURL, err := url.Parse("https://api.openweathermap.org/data/2.5/onecall")
	if err != nil {
		return WeatherAPIResponse{}, err
	}

	q := apiURL.Query()
	q.Add("lat", fmt.Sprint(lat))
	q.Add("lon", fmt.Sprint(lon))
	q.Add("exclude", "minutely,hourly,alerts")
	q.Add("units", "imperial")
	q.Add("appid", os.Getenv("WEATHER_API_KEY"))
	apiURL.RawQuery = q.Encode()

	httpResponse, err := http.Get(apiURL.String())
	if err != nil {
		return WeatherAPIResponse{}, err
	}

	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return WeatherAPIResponse{}, err
	}

	var response WeatherAPIResponse
	err = json.Unmarshal(responseBytes, &response)

	if err != nil {
		return WeatherAPIResponse{}, err
	}

	return response, err
}
