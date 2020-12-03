package ports

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/cesarFuhr/gocrypto/internal/domain/keys"
	"github.com/cesarFuhr/gocrypto/internal/ports/presenters"
)

type keyOpts struct {
	Scope      string `json:"scope"`
	Expiration string `json:"expiration"`
}

type keyHandler struct {
	service keys.KeyService
}

// KeyHandlerInterface describes a http handler interface
type KeyHandlerInterface interface {
	Post(w http.ResponseWriter, r *http.Request)
}

// NewKeyHandler creates a new http key handler
func NewKeyHandler(s keys.KeyService) KeyHandlerInterface {
	return keyHandler{
		service: s,
	}
}

// Post http translator
func (h *keyHandler) Post(w http.ResponseWriter, r *http.Request) {
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

	key, err := h.service.CreateKey(o.Scope, exp)
	if err != nil {
		internalServerError(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(presenters.NewHTTPCreateKey(key))
	return
}
