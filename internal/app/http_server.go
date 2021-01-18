package server

import (
	"net/http"

	"github.com/cesarFuhr/gocrypto/internal/app/ports"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// HTTPLogger http server logger
type HTTPLogger interface {
	Info(string, ...zap.Field)
}

// NewHTTPServer creates a new http handler
func NewHTTPServer(
	l HTTPLogger,
	kH ports.KeyHandler,
	eH ports.EncryptHandler,
	dH ports.DecryptHandler,
) *http.Server {
	router := mux.NewRouter()
	logger := newLoggerMiddleware(l)

	router.Use(logger)

	router.
		HandleFunc("/keys", kH.Post).
		Methods(http.MethodPost)
	router.
		HandleFunc("/keys/{keyID}", kH.Get).
		Methods(http.MethodGet)

	router.
		HandleFunc("/encrypt", eH.Post).
		Methods(http.MethodPost)

	router.
		HandleFunc("/decrypt", dH.Post).
		Methods(http.MethodPost)

	return &http.Server{
		Handler: router,
	}
}
