package server

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func newLoggerMiddleware(log HTTPLogger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		logger := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			wraped := wrapResponseWriter(w)
			h.ServeHTTP(wraped, r)

			log.Info(
				"HTTP Request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", wraped.Status()),
				zap.Int64("rSize", r.ContentLength),
				zap.Duration("latency", time.Since(startTime)),
			)
		})
		return logger
	}
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}
