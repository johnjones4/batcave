package main

import (
	"os"

	"github.com/johnjones4/hal-9000/server/hal9000/cli"
)

func main() {
	host := ""
	scheme := "https"
	if len(os.Args) > 1 {
		host = os.Args[1]
	} else {
		host = os.Getenv("HAL9000_HOST")
	}
	if len(os.Args) > 2 {
		scheme = os.Args[2]
	}
	c := cli.New(scheme, host, os.Getenv("HAL9000_SETTINGS_PATH"))
	c.Run()
}
