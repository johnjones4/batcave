package main

import (
	"fmt"
	"hal9000"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Booting up ...")
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = hal9000.BootUp()
	if err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/ws", wsHandler)

	fmt.Println("Ready")

	err = http.ListenAndServe(os.Getenv("HTTP_SERVER"), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
