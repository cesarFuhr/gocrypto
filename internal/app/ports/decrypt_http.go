package ports

import (
	"net/http"

	"github.com/cesarFuhr/gocrypto/internal/app/domain/crypto"
	"github.com/cesarFuhr/gocrypto/internal/app/domain/keys"
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
			replyJSON(w, http.StatusPreconditionFailed, HTTPError{
				Message: "Key was not found",
			})
			return
		}
		internalServerError(w)
		return
	}

	replyJSON(w, http.StatusOK, HTTPDecrypt{
		Data: string(decrypted),
	})
	return
}
