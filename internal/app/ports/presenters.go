package ports

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/cesarFuhr/gocrypto/internal/app/domain/keys"
)

// HTTPError Exception formatter to all http badRequests
type HTTPError struct {
	Message string `json:"message"`
}

// HTTPCreateKey Http representation of the create key response body
type HTTPCreateKey struct {
	KeyID      string `json:"keyID"`
	Expiration string `json:"expiration"`
	PublicKey  string `json:"publicKey"`
}

// NewHTTPCreateKey Builder for the http CreateKey response
func NewHTTPCreateKey(k keys.Key) HTTPCreateKey {
	return HTTPCreateKey{
		KeyID:      k.ID,
		Expiration: k.Expiration.UTC().Format(time.RFC3339),
		PublicKey:  formatPublicKey(k.Pub),
	}
}

func formatPublicKey(pubKey *rsa.PublicKey) string {
	b := base64.RawStdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(pubKey))
	return b
}

// HTTPEncrypt representation of the encrypt response body
type HTTPEncrypt struct {
	EncryptedData string `json:"encryptedData"`
}

// HTTPDecrypt representation of the encrypt response body
type HTTPDecrypt struct {
	Data string `json:"data"`
}

func replyJSON(w http.ResponseWriter, code int, obj interface{}) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(obj)
}

func internalServerError(w http.ResponseWriter) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(HTTPError{
		Message: "There was an unexpected error",
	})
}
