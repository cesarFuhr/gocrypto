package ports

import (
	"encoding/json"
	"net/http"

	"github.com/cesarFuhr/gocrypto/internal/domain/crypto"
	"github.com/cesarFuhr/gocrypto/internal/domain/keys"
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
	json.NewEncoder(w).Encode(HTTPEncrypt{
		EncryptedData: string(encrypted),
	})
	return
}
