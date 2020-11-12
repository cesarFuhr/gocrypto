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

type decrypt struct {
	KeyID         string `json:"keyID"`
	EncryptedData string `json:"encryptedData"`
}

type keyStoreInterface interface {
	CreateKey(string, time.Time) (keys.Key, error)
	FindKey(string) (keys.Key, error)
}

type cryptoInterface interface {
	Encrypt(*rsa.PublicKey, string) ([]byte, error)
	Decrypt(*rsa.PrivateKey, string) ([]byte, error)
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
	router.Handle("/decrypt", http.HandlerFunc(s.decryptHandler))

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
			json.NewEncoder(w).Encode(presenters.HTTPError{
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
	json.NewEncoder(w).Encode(presenters.HTTPEncrypt{
		EncryptedData: string(encrypted),
	})
	return
}

func (s *KeyServer) decryptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}

	var o decrypt
	decodeJSONBody(r, &o)

	key, err := s.keyStore.FindKey(o.KeyID)
	if err != nil {
		if err == keys.ErrKeyNotFound {
			w.WriteHeader(http.StatusPreconditionFailed)
			json.NewEncoder(w).Encode(presenters.HTTPError{
				Message: "Key was not found",
			})
			return
		}
		internalServerError(w)
		return
	}

	decrypted, err := s.crypto.Decrypt(key.Priv, o.EncryptedData)
	if err != nil {
		internalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(presenters.HTTPDecrypt{
		Data: string(decrypted),
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
			json.NewEncoder(w).Encode(presenters.HTTPError{
				Message: mr.msg,
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(presenters.HTTPError{
			Message: fmt.Sprint(err),
		})
		return
	}

	exp, err := time.Parse(time.RFC3339, o.Expiration)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(presenters.HTTPError{
			Message: "Invalid: expiration property format",
		})
		return
	}

	key, err := s.keyStore.CreateKey(o.Scope, exp)
	if err != nil {
		internalServerError(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(presenters.NewHTTPCreateKey(key))
	return
}

func (s *KeyServer) getKeys(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("keyID")

	key, err := s.keyStore.FindKey(id)
	if err != nil {
		if err == keys.ErrKeyNotFound {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(presenters.HTTPError{
				Message: "Key was not found",
			})
			return
		}
		internalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(presenters.NewHTTPCreateKey(key))
	return
}
