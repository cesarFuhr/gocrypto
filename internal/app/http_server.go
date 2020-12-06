package server

import (
	"net/http"

	"github.com/cesarFuhr/gocrypto/internal/app/ports"
)

// HTTPServer http server interface
type HTTPServer interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// HTTPLogger http server logger
type HTTPLogger interface {
	Info(...interface{})
}

type httpServer struct {
	log            HTTPLogger
	keysHandler    ports.KeyHandler
	encryptHandler ports.EncryptHandler
	decryptHandler ports.DecryptHandler
}

// NewHTTPServer creates a new http handler
func NewHTTPServer(
	l HTTPLogger,
	kH ports.KeyHandler,
	eH ports.EncryptHandler,
	dH ports.DecryptHandler,
) HTTPServer {
	return &httpServer{
		log:            l,
		keysHandler:    kH,
		encryptHandler: eH,
		decryptHandler: dH,
	}
}

func (s *httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := http.NewServeMux()
	logger := newLoggerMiddleware(s.log)

	router.Handle("/keys", logger(s.handleKeys(w, r)))
	router.Handle("/encrypt", logger(s.handleEncrypt(w, r)))
	router.Handle("/decrypt", logger(s.handleDecrypt(w, r)))

	router.ServeHTTP(w, r)
}

func (s *httpServer) handleKeys(w http.ResponseWriter, r *http.Request) http.Handler {
	f := methodNotAllowed

	switch r.Method {
	case http.MethodPost:
		f = s.keysHandler.Post
	case http.MethodGet:
		f = s.keysHandler.Get
	}

	return http.HandlerFunc(f)
}

func (s *httpServer) handleEncrypt(w http.ResponseWriter, r *http.Request) http.Handler {
	f := methodNotAllowed

	if r.Method == http.MethodPost {
		f = s.encryptHandler.Post
	}

	return http.HandlerFunc(f)
}

func (s *httpServer) handleDecrypt(w http.ResponseWriter, r *http.Request) http.Handler {
	f := methodNotAllowed

	if r.Method == http.MethodPost {
		f = s.decryptHandler.Post
	}

	return http.HandlerFunc(f)
}

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte{})
}
