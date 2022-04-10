package main

import (
	"os"
	"strconv"

	"github.com/johnjones4/hal-9000/server/hal9000/api"
	"github.com/johnjones4/hal-9000/server/hal9000/core"
	"github.com/johnjones4/hal-9000/server/hal9000/runtime"
	"github.com/johnjones4/hal-9000/server/hal9000/socket"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	runtime, err := runtime.New()
	if err != nil {
		panic(err)
	}

	defaulLat, err := strconv.ParseFloat(os.Getenv("DEFAULT_LATITUDE"), 64)
	if err != nil {
		panic(err)
	}

	defaulLon, err := strconv.ParseFloat(os.Getenv("DEFAULT_LONGITUDE"), 64)
	if err != nil {
		panic(err)
	}

	socket := &socket.Server{
		Host:    os.Getenv("SOCKET_HOST"),
		Runtime: runtime,
		Location: core.Coordinate{
			Latitude:  defaulLat,
			Longitude: defaulLon,
		},
	}
	startSocket := func() {
		err := socket.Run()
		panic(err)
	}
	go startSocket()

	api := &api.API{
		Host:    os.Getenv("HTTP_HOST"),
		Runtime: runtime,
	}
	err = api.Run()
	panic(err)
}
