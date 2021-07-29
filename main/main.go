package main

import (
	"encoding/json"
	"fmt"
	"hal9000"
	"os"

	"github.com/joho/godotenv"
)

func main() {
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

	ses, err := hal9000.InitiateNewSession(nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	response, err := ses.ProcessIncomingMessage(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(bytes))
}
