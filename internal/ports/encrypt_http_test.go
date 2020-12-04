package ports

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cesarFuhr/gocrypto/internal/domain/keys"
	"github.com/google/uuid"
)

type CryptoServiceStub struct {
	CalledWith []interface{}
}

func (s *CryptoServiceStub) Encrypt(keyID string, m string) ([]byte, error) {
	s.CalledWith = []interface{}{keyID, m}
	if m == "error" {
		return []byte{}, errors.New("some error")
	}
	if keyID == "notFound" {
		return []byte{}, keys.ErrKeyNotFound
	}
	return []byte{10, 10, 10}, nil
}

func TestEncrypt(t *testing.T) {
	cryptoStub := CryptoServiceStub{}
	h := NewEncryptHandler(&cryptoStub)
	t.Run("Should return a 200 if it was a success", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/encrypt", nil)
		response := httptest.NewRecorder()

		want := http.StatusOK

		h.Post(response, request)

		assertStatus(t, response.Code, want)
	})
	t.Run("Should return the correct properties", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/encrypt", nil)
		response := httptest.NewRecorder()

		h.Post(response, request)

		var got HTTPEncrypt
		json.NewDecoder(response.Body).Decode(&got)

		if got.EncryptedData == "" {
			t.Errorf("Expecting data prop, got %v", got)
		}
	})
	t.Run("Should call Encrypt with the right params", func(t *testing.T) {
		keyID := uuid.New().String()
		data := "testing"
		requestBody, _ := json.Marshal(map[string]string{
			"keyID": keyID,
			"data":  data,
		})

		request, _ := http.NewRequest(http.MethodPost, "/encrypt", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()
		h.Post(response, request)

		assertInsideSlice(t, cryptoStub.CalledWith, data)
		assertInsideSlice(t, cryptoStub.CalledWith, keyID)
	})
	t.Run("Should return a internal server error if there was a problem encrypting", func(t *testing.T) {
		keyID := uuid.New().String()
		data := "error"
		requestBody, _ := json.Marshal(map[string]string{
			"keyID": keyID,
			"data":  data,
		})

		request, _ := http.NewRequest(http.MethodPost, "/encrypt", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()
		h.Post(response, request)

		assertStatus(t, response.Code, http.StatusInternalServerError)
		assertInsideJSON(t, response.Body, "message", "There was an unexpected error")
	})
	t.Run("Should return a precondition fail if the key does not exists", func(t *testing.T) {
		requestBody, _ := json.Marshal(map[string]string{
			"keyID": "notFound",
		})
		request, _ := http.NewRequest(http.MethodPost, "/encrypt", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()
		h.Post(response, request)

		assertStatus(t, response.Code, http.StatusPreconditionFailed)
		assertInsideJSON(t, response.Body, "message", "Key was not found")
	})
}
