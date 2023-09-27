package api

import (
	"main/core"
	"net/http"
)

type conversationResponse struct {
	Type     string         `json:"type"`
	Request  *core.Request  `json:"request"`
	Response *core.Response `json:"response"`
}

func (a *API) converse(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		a.handleError(w, err, http.StatusBadRequest)
		return
	}
	defer c.Close()

	for {
		var req core.Request
		err = c.ReadJSON(&req)
		if err != nil {
			a.Log.Error(err)
			return
		}

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
