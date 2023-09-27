package api

import (
	"encoding/json"
	"main/core"
	"net/http"
)

func (a *API) authMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		clientId := r.Header.Get("X-Client-Id")
		if clientId == "" {
			a.handleError(w, nil, http.StatusUnauthorized)
			return
		}

		client, err := a.ClientRegistry.Client(r.Context(), "api", clientId, func(client *core.Client, info string) error {
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

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)

}
