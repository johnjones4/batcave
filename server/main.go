package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"main/services"
	"main/types"

	"github.com/gorilla/mux"
)

func loadConfig() (types.Configuration, error) {
	configFileBytes, err := os.ReadFile(os.Getenv("CONFIG_FILE_PATH"))
	if err != nil {
		return types.Configuration{}, err
	}

	var config types.Configuration
	err = json.Unmarshal(configFileBytes, &config)
	if err != nil {
		return types.Configuration{}, err
	}

	return config, nil
}

func start() error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	responseMutex := sync.Mutex{}
	errorChannel := make(chan error)
	newsChannel := make(chan []types.NewsItem)
	weatherChannel := make(chan []types.Weather)

	go services.StartNewsUpdater(config.RSSFeeds, newsChannel, errorChannel)
	go services.StartWeatherUpdater(config.WeatherLocations, weatherChannel, errorChannel)

	response := types.Response{
		IFrames: config.IFrames,
	}

	go func() {
		for {
			select {
			case err := <-errorChannel:
				log.Println(err)
			case newsItems := <-newsChannel:
				responseMutex.Lock()
				response.NewsItems = newsItems
				responseMutex.Unlock()
			case weather := <-weatherChannel:
				responseMutex.Lock()
				response.Weather = weather
				responseMutex.Unlock()
			}
		}
	}()

	r := mux.NewRouter()

	r.HandleFunc("/api/data", func(w http.ResponseWriter, req *http.Request) {
		responseMutex.Lock()
		jsonResponse(w, response)
		responseMutex.Unlock()
	}).Methods("GET")

	handler := logRequestHandler(r)

	srv := &http.Server{
		Addr:    os.Getenv("HTTP_SERVER"),
		Handler: handler,
	}
	return srv.ListenAndServe()
}

func main() {
	err := start()
	if err != nil {
		panic(err)
	}
}
