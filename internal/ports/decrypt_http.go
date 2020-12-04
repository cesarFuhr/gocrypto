package ports

import (
	"encoding/json"
	"net/http"

	"github.com/cesarFuhr/gocrypto/internal/domain/crypto"
	"github.com/cesarFuhr/gocrypto/internal/domain/keys"
)

type decryptReqBody struct {
	KeyID         string `json:"keyID"`
	EncryptedData string `json:"encryptedData"`
}

// DecryptHandler decrypt http handler
type DecryptHandler interface {
	Post(w http.ResponseWriter, r *http.Request)
}

type decryptHandler struct {
	service crypto.Service
}

// NewDecryptHandler creates a decrypt http handler
func NewDecryptHandler(s crypto.Service) DecryptHandler {
	return &decryptHandler{
		service: s,
	}
}

func (s *decryptHandler) Post(w http.ResponseWriter, r *http.Request) {
	var o decryptReqBody
	decodeJSONBody(r, &o)

	decrypted, err := s.service.Decrypt(o.KeyID, o.EncryptedData)
	if err != nil {
		if err == keys.ErrKeyNotFound {
			w.WriteHeader(http.StatusPreconditionFailed)
			json.NewEncoder(w).Encode(HTTPError{
				Message: "Key was not found",
			})
			return
		}
		internalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(HTTPDecrypt{
		Data: string(decrypted),
	})
	return
}
