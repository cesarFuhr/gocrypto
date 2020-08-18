package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"time"

	"github.com/cesarFuhr/gocrypto/keys"
)

type createKeyResponseBody struct {
	KeyID      string    `json:"keyID"`
	Expiration time.Time `json:"expiration"`
	PublicKey  string    `json:"publicKey"`
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
	key := s.keyStore.CreateKey("test", time.Now())

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createKeyResponseBody{
		KeyID:      key.ID,
		Expiration: key.Expiration.UTC(),
		PublicKey:  formatPublicKey(key.Pub),
	})
}

func formatPublicKey(pubKey *rsa.PublicKey) string {
	b := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(pubKey)})
	return string(b)
}
