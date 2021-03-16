package ports

import (
	"errors"
	"net/http"
	"time"

	"github.com/cesarFuhr/gocrypto/internal/app/domain/keys"
	"github.com/gorilla/mux"
)

type keyOpts struct {
	Scope      string `json:"scope" validate:"required,gt=0,lte=50"`
	Expiration string `json:"expiration" validate:"required,datetime"`
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
			replyJSON(w, mr.status, HTTPError{
				Message: mr.msg,
			})
			return
		}
		replyJSON(w, http.StatusInternalServerError, HTTPError{
			Message: err.Error(),
		})
	}

	exp, err := time.Parse(time.RFC3339, o.Expiration)
	if err != nil {
		replyJSON(w, http.StatusBadRequest, HTTPError{
			Message: "Invalid: expiration property format",
		})
		return
	}

	key, err := h.service.CreateKey(o.Scope, exp)
	if err != nil {
		internalServerError(w)
		return
	}

	replyJSON(w, http.StatusCreated, NewHTTPCreateKey(key))
}

func (h *keyHandler) Get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, ok := params["keyID"]

	if !ok {
		replyJSON(w, http.StatusBadRequest, HTTPError{
			Message: "missing path param keyID",
		})
		return
	}

	key, err := h.service.FindKey(id)
	if err != nil {
		if err == keys.ErrKeyNotFound {
			replyJSON(w, http.StatusNotFound, HTTPError{
				Message: "Key was not found",
			})
			return
		}
		internalServerError(w)
		return
	}

	replyJSON(w, http.StatusOK, NewHTTPCreateKey(key))
}
