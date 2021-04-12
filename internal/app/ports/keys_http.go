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
	service   keys.KeyService
	validator keysValidator
}

// KeyHandler describes a http handler interface
type KeyHandler interface {
	Post(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Find(w http.ResponseWriter, r *http.Request)
}

// NewKeyHandler creates a new http key handler
func NewKeyHandler(s keys.KeyService) KeyHandler {
	return &keyHandler{
		service:   s,
		validator: keysValidator{},
	}
}

// Post http translator
func (h *keyHandler) Post(w http.ResponseWriter, r *http.Request) {
	var o keyOpts
	if err := decodeJSONBody(r, &o); err != nil {
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
		return
	}

	if err := h.validator.PostValidator(o); err != nil {
		replyJSON(w, http.StatusBadRequest, HTTPError{
			Message: err.Error(),
		})
		return
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
	id := params["keyID"]

	if err := h.validator.GetValidator(id); err != nil {
		replyJSON(w, http.StatusBadRequest, HTTPError{
			Message: err.Error(),
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

func (h *keyHandler) Find(w http.ResponseWriter, r *http.Request) {
	scope := r.URL.Query().Get("scope")

	if err := h.validator.FindValidator(scope); err != nil {
		replyJSON(w, http.StatusBadRequest, HTTPError{
			Message: err.Error(),
		})
		return
	}

	key, err := h.service.FindKeysByScope(scope)
	if err != nil {
		internalServerError(w)
		return
	}

	replyJSON(w, http.StatusOK, NewHTTPFindKeys(key))
}
