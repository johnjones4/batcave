package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func (a *apiConcrete) handleError(w http.ResponseWriter, err error, status int) {
	if err == nil {
		err = errors.New(http.StatusText(status))
	}
	a.Log.Error(err)
	http.Error(w, http.StatusText(status), status)
}

func (a *apiConcrete) jsonResponse(w http.ResponseWriter, j any) {
	bytes, err := json.Marshal(j)
	if err != nil {
		a.handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.Write(bytes)
}

func (a *apiConcrete) readJson(req *http.Request, i any) error {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	a.Log.Debug(string(body))

	return json.Unmarshal(body, i)
}
