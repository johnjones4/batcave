package api

import (
	"encoding/json"
	"io"
	"main/core"
	"net/http"
)

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

	res, err := a.coreHandler(r.Context(), req)
	if err != nil {
		a.handleError(w, err, http.StatusInternalServerError)
		return
	}

	a.jsonResponse(w, res)
}
