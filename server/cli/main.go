package main

import (
	"os"

	"github.com/johnjones4/hal-9000/hal9000/cli"
)

func main() {
	c := cli.New(os.Args[1])
	c.Run()
}
