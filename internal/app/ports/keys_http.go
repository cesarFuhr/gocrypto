package ports

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/cesarFuhr/gocrypto/internal/app/domain/keys"
)

type keyOpts struct {
	Scope      string `json:"scope"`
	Expiration string `json:"expiration"`
}

type keyHandler struct {
	service keys.KeyService
}

// KeyHandler describes a http handler interface
type KeyHandler interface {
	Post(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
}

// NewKeyHandler creates a new http key handler
func NewKeyHandler(s keys.KeyService) KeyHandler {
	return &keyHandler{
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
			json.NewEncoder(w).Encode(HTTPError{
				Message: mr.msg,
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(HTTPError{
			Message: fmt.Sprint(err),
		})
		return
	}

	exp, err := time.Parse(time.RFC3339, o.Expiration)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(HTTPError{
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
	json.NewEncoder(w).Encode(NewHTTPCreateKey(key))
	return
}

func (h *keyHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("keyID")

	key, err := h.service.FindKey(id)
	if err != nil {
		if err == keys.ErrKeyNotFound {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(HTTPError{
				Message: "Key was not found",
			})
			return
		}
		internalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewHTTPCreateKey(key))
	return
}
