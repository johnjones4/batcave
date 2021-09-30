package main

import (
	"net/http"

	"main/handlers"
	"main/types"

	"github.com/gorilla/mux"
)

func makeConfigRoute(config *types.Configuration) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		jsonResponse(w, *config)
	}
}

func makeWeatherRoute(config *types.Configuration) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		weatherInfo := make([]types.Weather, len(config.WeatherLocations))
		for i, coord := range config.WeatherLocations {
			info, err := handlers.GetWeather(coord)
			if err != nil {
				errorResponse(w, err)
				return
			}
			weatherInfo[i] = info
		}
		jsonResponse(w, weatherInfo)
	}
}

func InitRouter(config *types.Configuration) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/api/configuration", makeConfigRoute(config)).Methods("GET")
	r.HandleFunc("/api/weather", makeWeatherRoute(config)).Methods("GET")

	var handler http.Handler = r
	return logRequestHandler(handler)
}
