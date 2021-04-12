package ports

import (
	"net/http"

	"github.com/cesarFuhr/gocrypto/internal/app/domain/keys"
)

type encryptReqBody struct {
	KeyID string `json:"keyID"`
	Data  string `json:"data"`
}

type EncryptionService interface {
	Encrypt(string, string) ([]byte, error)
}

type EncryptHandler struct {
	service   EncryptionService
	validator encryptValidator
}

// NewEncryptHandler creates a new http key handler
func NewEncryptHandler(s EncryptionService) EncryptHandler {
	return EncryptHandler{
		service:   s,
		validator: encryptValidator{},
	}
}

func (h *EncryptHandler) Post(w http.ResponseWriter, r *http.Request) {
	var o encryptReqBody
	decodeJSONBody(r, &o)

	if err := h.validator.PostValidator(o); err != nil {
		replyJSON(w, http.StatusBadRequest, HTTPError{
			Message: err.Error(),
		})
		return
	}

	encrypted, err := h.service.Encrypt(o.KeyID, o.Data)
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

	replyJSON(w, http.StatusOK, HTTPEncrypt{
		EncryptedData: string(encrypted),
	})
}
