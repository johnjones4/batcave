package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/johnjones4/hal-9000/server/hal9000/storage"
)

func fail(w http.ResponseWriter, err error) {
	log.Println(err)
	w.WriteHeader(http.StatusForbidden)
}

var requiredHeaders = []string{
	"User-Agent",
	"X-Request-Time",
}

func makeRequestVerifier(clientStore *storage.ClientStore) func(http.Handler) http.Handler {
	return func(upstream http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientId := r.Header.Get("User-Agent")
			signature := r.Header.Get("X-Signature")

			if clientId == "" || signature == "" {
				fail(w, errors.New("invalid header"))
				return
			}

			client, err := clientStore.GetClient(clientId)
			if err != nil {
				fail(w, errors.New("invalid Authorization header"))
				return
			}

			values := make([]string, len(requiredHeaders))
			for i, header := range requiredHeaders {
				value := r.Header.Get(header)
				if value == "" {
					fail(w, errors.New("missing required header"))
					return
				}
				values[i] = value
			}
			valueString := strings.Join(values, ":")

			h := hmac.New(sha256.New, []byte(client.Key))

			h.Write([]byte(valueString))

			sha := hex.EncodeToString(h.Sum(nil))

			if sha != signature {
				fail(w, errors.New("invalid signature"))
				return
			}

			reqTimeString := r.Header.Get("X-Request-Time")
			reqTime, err := time.Parse(time.RFC3339, reqTimeString)
			if err != nil {
				fail(w, err)
				return
			}

			now := time.Now()
			bounds := time.Second
			lower := now.Add(-bounds)
			upper := now.Add(bounds)
			if reqTime.Before(lower) || reqTime.After(upper) {
				fail(w, errors.New("bad request time"))
				return
			}

			upstream.ServeHTTP(w, r)
		})
	}
}
