package ports

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cesarFuhr/gocrypto/internal/app/domain/keys"
)

type DecryptionServiceStub struct {
	CalledWith []interface{}
}

func (s *DecryptionServiceStub) Decrypt(keyID string, m string) ([]byte, error) {
	s.CalledWith = []interface{}{keyID, m}
	if m == "error" {
		return []byte{}, errors.New("some error")
	}
	if m == "notFound" {
		return []byte{}, keys.ErrKeyNotFound
	}
	return []byte{10, 10, 10}, nil
}

func TestDecrypt(t *testing.T) {
	cryptoStub := DecryptionServiceStub{}
	h := NewDecryptHandler(&cryptoStub)
	t.Run("Should return a 200 if it was a success", func(t *testing.T) {
		requestBody, _ := json.Marshal(decryptReqBody{
			EncryptedData: "mensagem",
			KeyID:         "f6a4633a-65f5-42f8-a984-38d87e3513ee",
		})

		request, _ := http.NewRequest(http.MethodPost, "/decrypt", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()

		want := http.StatusOK

		h.Post(response, request)

		assertStatus(t, response.Code, want)
	})
	t.Run("Should return the correct properties", func(t *testing.T) {
		requestBody, _ := json.Marshal(decryptReqBody{
			EncryptedData: "mensagem",
			KeyID:         "f6a4633a-65f5-42f8-a984-38d87e3513ee",
		})

		request, _ := http.NewRequest(http.MethodPost, "/decrypt", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()

		h.Post(response, request)

		var got HTTPDecrypt
		json.NewDecoder(response.Body).Decode(&got)

		if got.Data == "" {
			t.Errorf("Expecting data prop, got %v", got.Data)
		}
	})
	t.Run("Should call Decrypt with the right params", func(t *testing.T) {
		requestBody, _ := json.Marshal(decryptReqBody{
			EncryptedData: "message",
			KeyID:         "f6a4633a-65f5-42f8-a984-38d87e3513ee",
		})
		request, _ := http.NewRequest(http.MethodPost, "/decrypt", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()
		h.Post(response, request)

		assertInsideSlice(t, cryptoStub.CalledWith, "message")
		assertInsideSlice(t, cryptoStub.CalledWith, "f6a4633a-65f5-42f8-a984-38d87e3513ee")
	})
	t.Run("Should return a internal server error if there was a problem decrypting", func(t *testing.T) {
		requestBody, _ := json.Marshal(decryptReqBody{
			EncryptedData: "error",
			KeyID:         "f6a4633a-65f5-42f8-a984-38d87e3513ee",
		})
		request, _ := http.NewRequest(http.MethodPost, "/decrypt", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()
		h.Post(response, request)

		assertStatus(t, response.Code, http.StatusInternalServerError)
		assertInsideJSON(t, response.Body, "message", "There was an unexpected error")
	})
	t.Run("Should return a precondition fail if the key does not exists", func(t *testing.T) {
		requestBody, _ := json.Marshal(map[string]string{
			"keyID":         "f6a4633a-65f5-42f8-a984-38d87e3513ee",
			"encryptedData": "notFound",
		})
		request, _ := http.NewRequest(http.MethodPost, "/decrypt", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()
		h.Post(response, request)

		assertStatus(t, response.Code, http.StatusPreconditionFailed)
		assertInsideJSON(t, response.Body, "message", "Key was not found")
	})
}
