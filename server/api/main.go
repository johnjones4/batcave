package main

import (
	"net/http"
	"os"

	"github.com/johnjones4/hal-9000/hal9000/api"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	userStoreFile := os.Getenv("USER_STORE_FILE")
	stateStoreFile := os.Getenv("STATE_STORE_FILE")
	logFile := os.Getenv("LOG_FILE")
	tokenKey := os.Getenv("TOKEN_KEY")

	handler, err := api.New(userStoreFile, stateStoreFile, logFile, tokenKey)
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(os.Getenv("HTTP_HOST"), handler)
	if err != nil {
		panic(err)
	}
}
