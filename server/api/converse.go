package api

import (
	"main/core"
	"net/http"
)

type conversationResponse struct {
	Type        string            `json:"type"`
	Request     *core.Request     `json:"request"`
	Response    *core.Response    `json:"response"`
	PushMessage *core.PushMessage `json:"push"`
}

func (a *API) converse(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil) //TODO make this an interface
	if err != nil {
		a.handleError(w, err, http.StatusBadRequest)
		return
	}
	defer c.Close()

	as := a.SocketSender.RegisterActiveSocket(r.Header.Get("X-Client-Id"), c.RemoteAddr().String())
	defer a.SocketSender.DeregisterActiveSocket(as.ClientId, as.ConnectionId)

	IncomingMessages := make(chan core.Request)
	go func() {
		for {
			var req core.Request
			err = c.ReadJSON(&req)
			if err != nil {
				a.Log.Error(err)
				return
			}
			IncomingMessages <- req
		}
	}()

	for {
		select {
		case push := <-as.Messages:
			err = c.WriteJSON(conversationResponse{
				Type:        "push",
				PushMessage: &push,
			})
			if err != nil {
				a.Log.Error(err)
				return
			}
		case req := <-IncomingMessages:
			req.Source = "api"

			err = a.prepareRequest(r.Context(), &req)
			if err != nil {
				a.Log.Error(err)
				return
			}

			err = c.WriteJSON(conversationResponse{
				Type:    "request",
				Request: &req,
			})
			if err != nil {
				a.Log.Error(err)
				return
			}

			resp, err := a.coreHandler(r.Context(), &req)
			if err != nil {
				a.Log.Error(err)
				return
			}

			err = c.WriteJSON(conversationResponse{
				Type:     "response",
				Response: &resp,
			})
			if err != nil {
				a.Log.Error(err)
				return
			}
		}
	}
}
