package main

import (
	"net/http"

	"main/handlers"
	"main/types"

	"github.com/gorilla/mux"
)

func makeDataRoute(config *types.Configuration) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		response := types.Response{
			IFrames: config.IFrames,
			Weather: make([]types.Weather, len(config.WeatherLocations)),
		}

		for i, coord := range config.WeatherLocations {
			info, err := handlers.GetWeather(coord)
			if err != nil {
				errorResponse(w, err)
				return
			}
			response.Weather[i] = info
		}

		jsonResponse(w, response)
	}
}

func InitRouter(config *types.Configuration) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/api/data", makeDataRoute(config)).Methods("GET")

	var handler http.Handler = r
	return logRequestHandler(handler)
}
