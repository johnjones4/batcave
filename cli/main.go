package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"hal9000"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var user = flag.String("user", "", "HAL 9000 user")
var ended = false

func inputCollector(collector chan<- string) {
	fmt.Print("HAL9000> ")
	reader := bufio.NewReader(os.Stdin)
	t, _ := reader.ReadString('\n')
	if ended {
		return
	}
	collector <- strings.TrimSpace(t)
}

func responseHandler(c *websocket.Conn, response chan<- string) {
	defer close(response)
	for {
		_, message, err := c.ReadMessage()
		if ended {
			return
		}
		if err != nil {
			fmt.Println("Error: ", err)
			ended = true
			return
		}
		var resp hal9000.ResponseMessage
		err = json.Unmarshal(message, &resp)
		if err != nil {
			fmt.Println("Error: ", err)
			ended = true
			return
		}
		response <- resp.Text
	}
}

func main() {
	flag.Parse()

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws", RawQuery: url.Values{"user": {*user}}.Encode()}
	fmt.Printf("connecting to %s\n", u.String())

	connection, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer connection.Close()

	inputChannel := make(chan string)
	go inputCollector(inputChannel)

	responseChannel := make(chan string)
	go responseHandler(connection, responseChannel)

	waitingForResponse := false
	for !ended {
		select {
		case input := <-inputChannel:
			waitingForResponse = true
			if strings.EqualFold("exit", input) {
				ended = true
				connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			req := hal9000.RequestMessage{Message: input}
			err := connection.WriteJSON(req)
			if err != nil {
				ended = true
				fmt.Println(err)
				connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
		case response := <-responseChannel:
			fmt.Println(response) //TODO break
			if waitingForResponse {
				waitingForResponse = false
				go inputCollector(inputChannel)
			}
		}
	}
}
