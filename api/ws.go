package main

import (
	"encoding/json"
	"fmt"
	"hal9000"
	"net/http"

	"github.com/gorilla/websocket"
)

type InterfaceTypeWebsocket struct {
	Connection *websocket.Conn
}

func (i InterfaceTypeWebsocket) Name() string {
	return "websocket"
}

func (i InterfaceTypeWebsocket) SendMessage(m hal9000.ResponseMessage) error {
	responseBytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = i.Connection.WriteMessage(websocket.TextMessage, responseBytes)
	if err != nil {
		return err
	}
	return nil
}

type HAL9000HTTPRequest struct {
	SessionID string                 `json:"sessionId"`
	Request   hal9000.RequestMessage `json:"request"`
}

type HAL9000HTTPResponse struct {
	SessionID string                  `json:"sessionId"`
	Response  hal9000.ResponseMessage `json:"response"`
}

var upgrader = websocket.Upgrader{}

func wsHandler(w http.ResponseWriter, req *http.Request) {
	c, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		errorResponse(w, err)
		return
	}

	defer c.Close()

	ses, err := hal9000.NewSession(InterfaceTypeWebsocket{c})
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		_, request, err := c.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		var halReq HAL9000HTTPRequest
		err = json.Unmarshal(request, &halReq)
		if err != nil {
			fmt.Println(err)
			return
		}

		response, err := ses.ProcessIncomingMessage(halReq.Request)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = ses.Interface.SendMessage(response)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = ses.Save()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
