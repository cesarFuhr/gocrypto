package server

import "net/http"

func newLoggerMiddleware(log HTTPLogger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		logger := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer log.Info(
				r.Method, " ",
				r.URL.Path, " ",
				"rSize:", r.ContentLength,
			)
			h.ServeHTTP(w, r)
		})
		return logger
	}
}
