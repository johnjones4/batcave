package api

import (
	"log"
	"main/intent"
	"main/learning"
	"main/storage"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func makeServer(intentSet *intent.IntentSet, userStore *storage.UserStore, stateStore *storage.StateStore, logger *learning.InteractionLogger, tm *security.TokenManager) http.Handler {
	r := mux.NewRouter()
	r.Use(
		logRequest,
	)
	r.Handle("/api/login", makeLoginHandler(userStore, tm)).Methods("POST")
	r.Handle("/api/request", makeRequestHandler(intentSet, userStore, stateStore, logger, tm)).Methods("POST")
	return r
}

type recordingResponseWriter struct {
	parent http.ResponseWriter
	status int
	length int
}

func (w *recordingResponseWriter) Header() http.Header {
	return w.parent.Header()
}

func (w *recordingResponseWriter) Write(b []byte) (int, error) {
	i, e := w.parent.Write(b)
	w.length += len(b)
	w.status = http.StatusOK
	return i, e
}

func (w *recordingResponseWriter) WriteHeader(statusCode int) {
	w.parent.WriteHeader(statusCode)
	w.status = statusCode
}

func logRequest(upstream http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := recordingResponseWriter{
			parent: w,
		}
		start := time.Now()
		upstream.ServeHTTP(&writer, r)
		span := time.Since(start)
		log.Printf("%s %s -> %d [%db] %s", r.Method, r.URL.String(), writer.status, writer.length, span.String())
	})
}

func handleError(w http.ResponseWriter, err error, code int) {
	log.Printf("[ERROR %s]", err.Error())
	w.WriteHeader(code)
}

func handleInternalError(w http.ResponseWriter, err error) {
	handleError(w, err, http.StatusInternalServerError)
}
