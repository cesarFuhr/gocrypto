package ports

import (
	"net/http"

	"github.com/cesarFuhr/gocrypto/internal/app/domain/crypto"
	"github.com/cesarFuhr/gocrypto/internal/app/domain/keys"
)

type encryptReqBody struct {
	KeyID string `json:"keyID"`
	Data  string `json:"data"`
}

// EncryptHandler encryption http handler
type EncryptHandler interface {
	Post(w http.ResponseWriter, r *http.Request)
}

type encryptHandler struct {
	service crypto.Service
}

// NewEncryptHandler creates a new http key handler
func NewEncryptHandler(s crypto.Service) EncryptHandler {
	return &encryptHandler{
		service: s,
	}
}

func (h *encryptHandler) Post(w http.ResponseWriter, r *http.Request) {
	var o encryptReqBody
	decodeJSONBody(r, &o)

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
