package main

import (
	"crypto/rsa"
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

type encrypt struct {
	KeyID string `json:"keyID"`
	Data  string `json:"data"`
}

type keyStoreInterface interface {
	CreateKey(string, time.Time) keys.Key
	FindKey(string) (keys.Key, error)
}

type cryptoInterface interface {
	Encrypt(*rsa.PublicKey, string) ([]byte, error)
}

// KeyServer key HTTP API server
type KeyServer struct {
	keyStore keyStoreInterface
	crypto   cryptoInterface
}

// ServeHTTP serves http requests
func (s *KeyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	router := http.NewServeMux()
	router.Handle("/keys", http.HandlerFunc(s.keysHandler))
	router.Handle("/encrypt", http.HandlerFunc(s.encryptHandler))

	router.ServeHTTP(w, r)
}

func (s *KeyServer) encryptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}

	var o encrypt
	decodeJSONBody(r, &o)

	key, err := s.keyStore.FindKey(o.KeyID)
	if err != nil {
		if err == keys.ErrKeyNotFound {
			w.WriteHeader(http.StatusPreconditionFailed)
			json.NewEncoder(w).Encode(presenters.HttpError{
				Message: "Key was not found",
			})
			return
		}
		internalServerError(w)
		return
	}

	encrypted, err := s.crypto.Encrypt(key.Pub, o.Data)
	if err != nil {
		internalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(presenters.HttpEncrypt{
		EncryptedData: string(encrypted),
	})
	return
}

func (s *KeyServer) keysHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getKeys(w, r)
	case http.MethodPost:
		s.createKeys(w, r)
	default:
		methodNotAllowed(w)
	}
	return
}

func (s *KeyServer) createKeys(w http.ResponseWriter, r *http.Request) {
	var o keyOpts
	err := decodeJSONBody(r, &o)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			w.WriteHeader(mr.status)
			json.NewEncoder(w).Encode(presenters.HttpError{
				Message: mr.msg,
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
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
	return
}

func (s *KeyServer) getKeys(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("keyID")

	key, err := s.keyStore.FindKey(id)
	if err != nil {
		if err == keys.ErrKeyNotFound {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(presenters.HttpError{
				Message: "Key was not found",
			})
			return
		}
		internalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(presenters.NewHttpCreateKey(key))
	return
}
