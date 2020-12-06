package ports

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/cesarFuhr/gocrypto/internal/app/domain/keys"
	"github.com/google/uuid"
)

type KeyServiceStub struct {
	CalledWith       []interface{}
	LastDeliveredKey keys.Key
}

var rsaKey, _ = rsa.GenerateKey(rand.Reader, 4098)

func (s *KeyServiceStub) CreateKey(scope string, exp time.Time) (keys.Key, error) {
	s.CalledWith = []interface{}{scope, exp}
	if scope == "ERROR" {
		return keys.Key{}, errors.New("A ERROR")
	}
	return keys.Key{
		Scope:      scope,
		Expiration: time.Now().AddDate(0, 0, 1),
		ID:         uuid.New().String(),
		Pub:        &rsaKey.PublicKey,
		Priv:       rsaKey,
	}, nil
}

func (s *KeyServiceStub) FindKey(id string) (keys.Key, error) {
	s.CalledWith = []interface{}{id}
	if id == "notFound" {
		return keys.Key{}, keys.ErrKeyNotFound
	}
	if id == "otherError" {
		return keys.Key{}, errors.New("Any error at all")
	}

	s.LastDeliveredKey = keys.Key{
		Scope:      "scope",
		Expiration: time.Now().AddDate(0, 0, 1),
		ID:         id,
		Pub:        &rsaKey.PublicKey,
		Priv:       rsaKey,
	}
	return s.LastDeliveredKey, nil
}

func (s *KeyServiceStub) FindScopedKey(id string, scope string) (keys.Key, error) {
	s.CalledWith = []interface{}{id, scope}
	if id == "notFound" {
		return keys.Key{}, keys.ErrKeyNotFound
	}
	if id == "otherError" {
		return keys.Key{}, errors.New("Any error at all")
	}

	s.LastDeliveredKey = keys.Key{
		Scope:      "scope",
		Expiration: time.Now().AddDate(0, 0, 1),
		ID:         id,
		Pub:        &rsaKey.PublicKey,
		Priv:       rsaKey,
	}
	return s.LastDeliveredKey, nil
}

var validReqBody, _ = json.Marshal(keyOpts{"scope", time.Now().UTC().Format(time.RFC3339)})

func TestPOSTKeys(t *testing.T) {
	keyServiceStub := KeyServiceStub{}
	h := NewKeyHandler(&keyServiceStub)
	t.Run("Should return 201 on /keys", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/keys", bytes.NewBuffer(validReqBody))
		response := httptest.NewRecorder()

		h.Post(response, request)

		assertStatus(t, response.Code, http.StatusCreated)
	})
	t.Run("Should return valid json on /keys", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/keys", bytes.NewBuffer(validReqBody))
		response := httptest.NewRecorder()

		h.Post(response, request)

		respBytes, _ := ioutil.ReadAll(response.Body)

		if !json.Valid(respBytes) {
			t.Errorf("got an invalid JSON %q", respBytes)
		}
	})
	t.Run("Should return all properties on /keys response", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/keys", bytes.NewBuffer(validReqBody))
		response := httptest.NewRecorder()

		wants := []string{"publicKey", "keyID", "expiration"}

		h.Post(response, request)
		respMap := map[string]interface{}{}
		extractJSON(response.Body, respMap)

		for _, want := range wants {
			if _, ok := respMap[want]; ok != true {
				t.Errorf("does not have the %q prop", want)
			}
		}
	})
	t.Run("Should call the CreateKey with expiration and scope", func(t *testing.T) {
		scope := "testing"
		expiration := time.Now().UTC().AddDate(0, 0, 1).Format(time.RFC3339)
		requestBody, _ := json.Marshal(map[string]string{
			"scope":      scope,
			"expiration": expiration,
		})
		request, _ := http.NewRequest(http.MethodPost, "/keys", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()

		h.Post(response, request)

		wantedExpiration, _ := time.Parse(time.RFC3339, expiration)
		assertInsideSlice(t, keyServiceStub.CalledWith, wantedExpiration)
		assertInsideSlice(t, keyServiceStub.CalledWith, scope)
	})
	t.Run("Should return a internal server error if there was an error creating keys", func(t *testing.T) {
		scope := "ERROR"
		expiration := time.Now().UTC().AddDate(0, 0, 1).Format(time.RFC3339)
		requestBody, _ := json.Marshal(map[string]string{
			"scope":      scope,
			"expiration": expiration,
		})
		request, _ := http.NewRequest(http.MethodPost, "/keys", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()

		h.Post(response, request)

		assertStatus(t, response.Code, http.StatusInternalServerError)
		assertInsideJSON(t, response.Body, "message", "There was an unexpected error")
	})
	t.Run("Should return a BadRequest if body is nil", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/keys", nil)
		response := httptest.NewRecorder()

		h.Post(response, request)

		assertStatus(t, response.Code, http.StatusBadRequest)
		assertInsideJSON(t, response.Body, "message", "Invalid: Empty body")
	})
	t.Run("Should return a BadRequest if expiration is not a RFC3339", func(t *testing.T) {
		scope := "testing"
		expiration := time.Now().UTC().AddDate(0, 0, 1).Format(time.RubyDate)
		requestBody, _ := json.Marshal(map[string]string{
			"scope":      scope,
			"expiration": expiration,
		})
		request, _ := http.NewRequest(http.MethodPost, "/keys", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()

		h.Post(response, request)

		assertStatus(t, response.Code, http.StatusBadRequest)
		assertInsideJSON(t, response.Body, "message", "Invalid: expiration property format")
	})
}

func TestGETKeys(t *testing.T) {
	keyServiceStub := KeyServiceStub{}
	h := NewKeyHandler(&keyServiceStub)
	t.Run("Should return a 200 if it was a success", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/keys", nil)
		response := httptest.NewRecorder()

		want := http.StatusOK

		h.Get(response, request)

		assertStatus(t, response.Code, want)
	})
	t.Run("Should return a Key if it was a success", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/keys", nil)
		response := httptest.NewRecorder()

		wants := []string{"publicKey", "keyID", "expiration"}

		h.Get(response, request)
		respMap := map[string]interface{}{}
		extractJSON(response.Body, respMap)

		for _, want := range wants {
			if _, ok := respMap[want]; ok != true {
				t.Errorf("does not have the %q prop", want)
			}
		}
	})
	t.Run("Should call find Keys with the right params", func(t *testing.T) {
		scope := "testing"
		expiration := time.Now().UTC().AddDate(0, 0, 1).Format(time.RFC3339)
		requestBody, _ := json.Marshal(map[string]string{
			"scope":      scope,
			"expiration": expiration,
		})

		request, _ := http.NewRequest(http.MethodPost, "/keys", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()
		h.Post(response, request)
		respMap := map[string]interface{}{}
		extractJSON(response.Body, respMap)

		getRequest, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/keys?keyID=%v", respMap["keyID"]), nil)
		h.Get(response, getRequest)

		assertInsideSlice(t, keyServiceStub.CalledWith, respMap["keyID"])
	})
	t.Run("If key was not found", func(t *testing.T) {
		t.Run("Should return a 404", func(t *testing.T) {
			want := http.StatusNotFound
			getRequest, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/keys?keyID=notFound"), nil)
			response := httptest.NewRecorder()
			h.Get(response, getRequest)

			assertStatus(t, response.Code, want)
			assertInsideJSON(t, response.Body, "message", "Key was not found")
		})
	})
	t.Run("If there was any other error", func(t *testing.T) {
		t.Run("Should return a 500", func(t *testing.T) {
			want := http.StatusInternalServerError
			getRequest, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/keys?keyID=otherError"), nil)
			response := httptest.NewRecorder()
			h.Get(response, getRequest)

			assertStatus(t, response.Code, want)
			assertInsideJSON(t, response.Body, "message", "There was an unexpected error")
		})
	})
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("want %d, got %d", want, got)
	}
}

func extractJSON(jBuff *bytes.Buffer, m map[string]interface{}) error {
	respBytes, _ := ioutil.ReadAll(jBuff)
	_ = json.Unmarshal(respBytes, &m)
	return nil
}

func assertInsideJSON(t *testing.T, jBuff *bytes.Buffer, wantedKey string, wantedValue interface{}) {
	t.Helper()
	got := map[string]interface{}{}
	extractJSON(jBuff, got)
	if !reflect.DeepEqual(got[wantedKey], wantedValue) {
		t.Errorf("got %v, want %v", got[wantedKey], wantedValue)
	}
}

func assertInsideSlice(t *testing.T, a []interface{}, want interface{}) {
	t.Helper()
	has := false
	for _, v := range a {
		if v == want {
			has = true
		}
	}
	if !has {
		t.Errorf("Did not found: %v, of type %T in %v", want, want, a)
	}
}
