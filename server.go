package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/cesarFuhr/gocrypto/keys"
	"github.com/cesarFuhr/gocrypto/presenters"
)

type keyOpts struct {
	Scope      string `json:"scope"`
	Expiration string `json:"expiration"`
}

type keyStoreInterface interface {
	CreateKey(string, time.Time) keys.Key
}

// KeyServer key HTTP API server
type KeyServer struct {
	keyStore keyStoreInterface
}

// ServeHTTP serves http requests
func (s *KeyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	router := http.NewServeMux()
	router.Handle("/keys", http.HandlerFunc(s.keysHandler))

	router.ServeHTTP(w, r)
}

func (s *KeyServer) keysHandler(w http.ResponseWriter, r *http.Request) {
	var o keyOpts
	err := decodeJSONBody(r, &o)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			w.WriteHeader(mr.status)
			json.NewEncoder(w).Encode(presenters.HttpError{
				Message: mr.msg,
			})
		}
		json.NewEncoder(w).Encode(presenters.HttpError{
			Message: fmt.Sprint(err),
		})
		return
	}

	exp, err := time.Parse(time.RFC3339, o.Expiration)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(presenters.HttpError{
			Message: "Invalid: expiration property format",
		})
		return
	}

	key := s.keyStore.CreateKey(o.Scope, exp)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(presenters.NewHttpCreateKey(key))
}
