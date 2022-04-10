package api

import (
	"log"
	"net/http"
	"time"
)

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
