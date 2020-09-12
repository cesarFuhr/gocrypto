package presenters

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"time"

	"github.com/cesarFuhr/gocrypto/keys"
)

// HttpBadRequest Exception formatter to all http badRequests
type HttpError struct {
	Message string `json:"message"`
}

// HttpCreateKey Http representation of the create key response body
type HttpCreateKey struct {
	KeyID      string `json:"keyID"`
	Expiration string `json:"expiration"`
	PublicKey  string `json:"publicKey"`
}

func NewHttpCreateKey(k keys.Key) HttpCreateKey {
	return HttpCreateKey{
		KeyID:      k.ID,
		Expiration: k.Expiration.UTC().Format(time.RFC3339),
		PublicKey:  formatPublicKey(k.Pub),
	}
}

func formatPublicKey(pubKey *rsa.PublicKey) string {
	b := pem.EncodeToMemory(&pem.Block{
		Type: "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(pubKey),
	})
	return string(b)
}
