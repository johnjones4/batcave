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

func (a *apiConcrete) directHandler(w http.ResponseWriter, r *http.Request) {
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

	if req.EventId == "" || req.ClientID == "" {
		a.handleError(w, nil, http.StatusBadRequest)
		return
	}

	client, err := a.ClientRegistry.Client(r.Context(), "api", req.ClientID, func(client *core.Client, info string) error {
		var receiver directAPIClientInfo
		err := json.Unmarshal([]byte(info), &receiver)
		if err != nil {
			return err
		}
		client.Info = receiver
		return nil
	})
	if err != nil {
		a.handleError(w, err, http.StatusInternalServerError)
		return
	}
	if client.Id == "" {
		a.handleError(w, nil, http.StatusUnauthorized)
		return
	}

	clientInfo, ok := client.Info.(directAPIClientInfo)
	if !ok {
		a.handleError(w, nil, http.StatusUnauthorized)
		return
	}

	headerKey := r.Header.Get("X-Api-Key")
	if headerKey == "" || clientInfo.Key != headerKey {
		a.handleError(w, nil, http.StatusUnauthorized)
		return
	}

	res, err := a.coreHandler(r.Context(), req)
	if err != nil {
		a.handleError(w, err, http.StatusInternalServerError)
		return
	}

	a.jsonResponse(w, res)
}
