package main

import (
	"os"

	"github.com/johnjones4/hal-9000/server/hal9000/api"
	"github.com/johnjones4/hal-9000/server/hal9000/runtime"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	runtime, err := runtime.New()
	if err != nil {
		panic(err)
	}

	api := &api.API{
		Host:    os.Getenv("HTTP_HOST"),
		Runtime: runtime,
	}
	err = api.Run()
	panic(err)
}
