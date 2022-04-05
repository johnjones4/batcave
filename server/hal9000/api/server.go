package api

import (
	"log"
	"net/http"
	"time"

	"github.com/johnjones4/hal-9000/server/hal9000/intent"
	"github.com/johnjones4/hal-9000/server/hal9000/learning"
	"github.com/johnjones4/hal-9000/server/hal9000/security"
	"github.com/johnjones4/hal-9000/server/hal9000/storage"

	"github.com/swaggest/rest/response/gzip"
	"github.com/swaggest/rest/web"
)

func makeServer(intentSet *intent.IntentSet, userStore *storage.UserStore, stateStore *storage.StateStore, logger *learning.InteractionLogger, tm *security.TokenManager) http.Handler {
	s := web.DefaultService()

	// Init API documentation schema.
	s.OpenAPI.Info.Title = "Basic Example"
	s.OpenAPI.Info.WithDescription("This app showcases a trivial REST API.")
	s.OpenAPI.Info.Version = "v1.2.3"

	// Setup middlewares.
	s.Use(
		logRequest,
		gzip.Middleware,
	)

	s.Post("/api/login", makeLoginHandler(userStore, tm))
	s.Post("/api/request", makeRequestHandler(intentSet, userStore, stateStore, logger, tm))

	return s
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
	return i, e
}

func (w *recordingResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.parent.WriteHeader(statusCode)
}

func logRequest(upstream http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := &recordingResponseWriter{
			parent: w,
			status: http.StatusOK,
		}
		start := time.Now()
		upstream.ServeHTTP(writer, r)
		span := time.Since(start)
		log.Printf("%s %s -> %d [%db] %s", r.Method, r.URL.String(), writer.status, writer.length, span.String())
	})
}
