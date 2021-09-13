package main

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"hal9000"
	"hal9000/util"
	"net/http"

	"github.com/gorilla/websocket"
)

type InterfaceTypeWebsocket struct {
	Connection *websocket.Conn
	Open       bool
}

func (i InterfaceTypeWebsocket) Type() string {
	return "websocket"
}

func (i InterfaceTypeWebsocket) ID() string {
	h := sha1.New()
	h.Write([]byte(i.Connection.RemoteAddr().String()))
	bs := h.Sum(nil)
	return fmt.Sprintf("ws-%x", bs)
}

func (i InterfaceTypeWebsocket) IsStillValid() bool {
	return i.Open
}

func (i InterfaceTypeWebsocket) SendMessage(m util.ResponseMessage) error {
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

var upgrader = websocket.Upgrader{}

func wsHandler(w http.ResponseWriter, req *http.Request) {
	userId := req.URL.Query().Get("user")
	if userId == "" {
		errorResponse(w, errors.New("no user id provided"))
		return
	}

	person, err := hal9000.GetPersonByID(userId)
	if err != nil {
		errorResponse(w, err)
		return
	}

	c, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		errorResponse(w, err)
		return
	}

	defer c.Close()

	iface := InterfaceTypeWebsocket{c, true}
	hal9000.RegisterTransientInterface(person, iface)
	ses := hal9000.NewSession(person, iface)

	for {
		_, request, err := c.ReadMessage()
		if err != nil {
			fmt.Println(err)
			iface.Open = false
			return
		}

		var halReq hal9000.RequestMessage
		err = json.Unmarshal(request, &halReq)
		if err != nil {
			fmt.Println(err)
			iface.Open = false
			return
		}

		response, err := ses.ProcessIncomingMessage(halReq)
		if err != nil {
			fmt.Println(err)
			iface.Open = false
			return
		}

		err = ses.Interface.SendMessage(response)
		if err != nil {
			fmt.Println(err)
			iface.Open = false
			return
		}

		hal9000.SaveSession(ses)
		if err != nil {
			fmt.Println(err)
			iface.Open = false
			return
		}
	}
}
