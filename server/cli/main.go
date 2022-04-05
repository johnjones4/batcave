package main

import (
	"os"

	"github.com/johnjones4/hal-9000/server/hal9000/cli"
)

func main() {
	host := ""
	if len(os.Args) > 1 {
		host = os.Args[1]
	} else {
		host = os.Getenv("HAL9000_HOST")
	}
	c := cli.New(host, os.Getenv("HAL9000_TOKEN_PATH"))
	c.Run()
}
