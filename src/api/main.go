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
	godotenv.Load()

	runtime, err := hal9000.BootUp()
	if err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/ws", wsHandler(runtime))
	http.HandleFunc("/sms", handleSMS(runtime))

	fmt.Println("Ready")

	err = http.ListenAndServe(os.Getenv("HTTP_SERVER"), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
