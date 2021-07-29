package main

import (
	"bufio"
	"fmt"
	"hal9000"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func shouldContinue(text string) bool {
	if strings.EqualFold("exit", text) {
		return false
	}
	return true
}

func get(r *bufio.Reader) string {
	t, _ := r.ReadString('\n')
	return strings.TrimSpace(t)
}

func printRepl() {
	fmt.Print("HAL9000> ")
}

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

	ses, err := hal9000.NewSession(hal9000.InterfaceTypeTerminal{os.Stdout})
	if err != nil {
		fmt.Println(err)
		return
	}

	printRepl()
	reader := bufio.NewReader(os.Stdin)
	text := get(reader)
	for ; shouldContinue(text); text = get(reader) {
		response, err := ses.ProcessIncomingMessage(hal9000.RequestMessage{text})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(response.Text)
		printRepl()
	}

	// bytes, err := json.Marshal(response)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(string(bytes))
}
