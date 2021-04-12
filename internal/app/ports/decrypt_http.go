package ports

import (
	"net/http"

	"github.com/cesarFuhr/gocrypto/internal/app/domain/keys"
)

type decryptReqBody struct {
	KeyID         string `json:"keyID"`
	EncryptedData string `json:"encryptedData"`
}

type DecryptHandler struct {
	service   DecryptionService
	validator decryptValidator
}

type DecryptionService interface {
	Decrypt(string, string) ([]byte, error)
}

// NewDecryptHandler creates a decrypt http handler
func NewDecryptHandler(s DecryptionService) DecryptHandler {
	return DecryptHandler{
		validator: decryptValidator{},
		service:   s,
	}
}

func (s *DecryptHandler) Post(w http.ResponseWriter, r *http.Request) {
	var o decryptReqBody
	decodeJSONBody(r, &o)

	if err := s.validator.PostValidator(o); err != nil {
		replyJSON(w, http.StatusBadRequest, HTTPError{
			Message: err.Error(),
		})
		return
	}

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
}
