package main

import (
	"context"
	"main/runtime"
)

func main() {
	rt, err := runtime.New(context.Background())
	if err != nil {
		panic(err)
	}
	err = rt.Start()
	if err != nil {
		panic(err)
	}
}
