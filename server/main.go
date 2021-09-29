package main

import (
	"encoding/json"
	"net/http"
	"os"

	"main/types"
)

func start() error {
	configFileBytes, err := os.ReadFile(os.Getenv("CONFIG_FILE_PATH"))
	if err != nil {
		return err
	}

	var config types.Config
	err = json.Unmarshal(configFileBytes, &config)
	if err != nil {
		return err
	}

	handler := InitRouter(&config)
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
