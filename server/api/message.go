package api

import (
	"encoding/json"
	"io"
	"main/core"
	"net/http"
)

type directAPIClientInfo struct {
	Key string `json:"key"`
}

func (a *API) message(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.handleError(w, err, http.StatusInternalServerError)
		return
	}

	var req core.Request
	err = json.Unmarshal(body, &req)
	if err != nil {
		a.handleError(w, err, http.StatusBadRequest)
		return
	}

	req.Source = "api"

	if req.EventId == "" || req.ClientID == "" || req.ClientID != r.Header.Get("X-Client-Id") {
		a.handleError(w, nil, http.StatusBadRequest)
		return
	}

	err = a.prepareRequest(r.Context(), &req)
	if err != nil {
		a.handleError(w, err, http.StatusInternalServerError)
		return
	}

	res, err := a.coreHandler(r.Context(), &req)
	if err != nil {
		a.handleError(w, err, http.StatusInternalServerError)
		return
	}

	a.jsonResponse(w, res)
}
