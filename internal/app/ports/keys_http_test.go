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
	"strings"
	"testing"
	"time"

	"github.com/cesarFuhr/gocrypto/internal/app/domain/keys"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type KeyServiceStub struct {
	CalledWith       []interface{}
	LastDeliveredKey keys.Key
	nextError        error
	nextFindResult   []keys.Key
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
	if s.nextError != nil {
		return keys.Key{}, s.nextError
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

func (s *KeyServiceStub) SetErrorFindKey(err error) {
	s.nextError = err
}

func (s *KeyServiceStub) SetErrorFindKeyByScope(err error) {
	s.nextError = err
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

func (s *KeyServiceStub) FindKeysByScope(scope string) ([]keys.Key, error) {
	s.CalledWith = []interface{}{scope}
	if s.nextError != nil {
		return nil, s.nextError
	}
	if s.nextFindResult != nil {
		return s.nextFindResult, nil
	}

	s.LastDeliveredKey = keys.Key{
		Scope:      "scope",
		Expiration: time.Now().AddDate(0, 0, 1),
		ID:         uuid.NewString(),
		Pub:        &rsaKey.PublicKey,
		Priv:       rsaKey,
	}
	return []keys.Key{s.LastDeliveredKey}, nil
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
		assertErrorMessage(t, response.Body, "message", "expiration is invalid")
	})
}

func TestGETKeys(t *testing.T) {
	keyServiceStub := KeyServiceStub{}
	h := NewKeyHandler(&keyServiceStub)
	m := make(map[string]string)
	t.Run("Should return a 200 if it was a success", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/keys/f6a4633a-65f5-42f8-a984-38d87e3513ee", nil)
		m["keyID"] = "f6a4633a-65f5-42f8-a984-38d87e3513ee"
		response := httptest.NewRecorder()

		want := http.StatusOK

		h.Get(response, mux.SetURLVars(request, m))

		assertStatus(t, response.Code, want)
	})
	t.Run("Should return a Key if it was a success", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/keys/f6a4633a-65f5-42f8-a984-38d87e3513ee", nil)
		m["keyID"] = "f6a4633a-65f5-42f8-a984-38d87e3513ee"
		response := httptest.NewRecorder()

		wants := []string{"publicKey", "keyID", "expiration"}

		h.Get(response, mux.SetURLVars(request, m))
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

		getRequest, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/keys/%v", respMap["keyID"]), nil)
		m["keyID"] = fmt.Sprint(respMap["keyID"])
		h.Get(response, mux.SetURLVars(getRequest, m))

		assertInsideSlice(t, keyServiceStub.CalledWith, respMap["keyID"])
	})
	t.Run("If key was not found", func(t *testing.T) {
		t.Run("Should return a 404", func(t *testing.T) {
			want := http.StatusNotFound
			getRequest, _ := http.NewRequest(http.MethodGet, "/keys/f6a4633a-65f5-42f8-a984-38d87e3513ee", nil)
			m["keyID"] = "f6a4633a-65f5-42f8-a984-38d87e3513ee"
			keyServiceStub.SetErrorFindKey(keys.ErrKeyNotFound)
			response := httptest.NewRecorder()

			h.Get(response, mux.SetURLVars(getRequest, m))

			assertStatus(t, response.Code, want)
			assertInsideJSON(t, response.Body, "message", "Key was not found")
		})
	})
	t.Run("If there was any other error", func(t *testing.T) {
		t.Run("Should return a 500", func(t *testing.T) {
			want := http.StatusInternalServerError
			getRequest, _ := http.NewRequest(http.MethodGet, "/keys/f6a4633a-65f5-42f8-a984-38d87e3513ee", nil)
			m["keyID"] = "f6a4633a-65f5-42f8-a984-38d87e3513ee"
			keyServiceStub.SetErrorFindKey(errors.New("another error"))
			response := httptest.NewRecorder()
			h.Get(response, mux.SetURLVars(getRequest, m))

			assertStatus(t, response.Code, want)
			assertInsideJSON(t, response.Body, "message", "There was an unexpected error")
		})
	})
}

func TestFindKeys(t *testing.T) {
	keyServiceStub := KeyServiceStub{}
	h := NewKeyHandler(&keyServiceStub)
	t.Run("Should return a 200 if it was a success", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/keys?scope=scope", nil)
		response := httptest.NewRecorder()

		want := http.StatusOK

		h.Find(response, request)

		assertStatus(t, response.Code, want)
	})
	t.Run("Should return .a list of Key if it was a success", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/keys?scope=scope", nil)
		response := httptest.NewRecorder()

		wants := []string{"publicKey", "keyID", "expiration"}

		h.Find(response, request)
		respArr := []map[string]interface{}{}
		json.Unmarshal(response.Body.Bytes(), &respArr)

		for _, want := range wants {
			for _, e := range respArr {
				if _, ok := e[want]; ok != true {
					t.Errorf("does not have the %q prop", want)
				}
			}
		}
	})
	t.Run("Should call find Keys by scope with the right params", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/keys?scope=target", nil)
		response := httptest.NewRecorder()

		h.Find(response, request)

		assertInsideSlice(t, keyServiceStub.CalledWith, "target")
	})
	t.Run("If scope is not well formatted", func(t *testing.T) {
		want := http.StatusBadRequest

		request, _ := http.NewRequest(http.MethodGet, "/keys?scope=1231231313132313123213131231231313123131313123141414", nil)
		response := httptest.NewRecorder()

		h.Find(response, request)
		respErr := HTTPError{}
		json.Unmarshal(response.Body.Bytes(), &respErr)

		assertStatus(t, response.Code, want)
		if !strings.Contains(respErr.Message, "scope is invalid") {
			t.Errorf("got %s, want %s", respErr, "some scope error message")
		}
	})
	t.Run("If no key was found", func(t *testing.T) {
		t.Run("Should return a 200 with no keys", func(t *testing.T) {
			want := http.StatusOK
			getRequest, _ := http.NewRequest(http.MethodGet, "/keys?scope=notFound", nil)
			response := httptest.NewRecorder()

			keyServiceStub.nextFindResult = []keys.Key{}

			h.Find(response, getRequest)
			respArr := []map[string]interface{}{}
			json.Unmarshal(response.Body.Bytes(), &respArr)

			assertStatus(t, response.Code, want)
			if len(respArr) != 0 {
				t.Errorf("got %d, want %d", len(respArr), 0)
			}
		})
	})
	t.Run("If there was any other error", func(t *testing.T) {
		t.Run("Should return a 500", func(t *testing.T) {
			want := http.StatusInternalServerError
			getRequest, _ := http.NewRequest(http.MethodGet, "/keys?scope=scope", nil)
			keyServiceStub.SetErrorFindKey(errors.New("another error"))
			response := httptest.NewRecorder()

			h.Find(response, getRequest)

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

func assertErrorMessage(t *testing.T, jBuff *bytes.Buffer, key, toMatch string) {
	t.Helper()
	json := map[string]interface{}{}
	extractJSON(jBuff, json)
	message := json[key]
	got := reflect.ValueOf(message)

	if !strings.Contains(got.String(), toMatch) {
		t.Errorf("got %v, want %v", got, toMatch)
	}
}
