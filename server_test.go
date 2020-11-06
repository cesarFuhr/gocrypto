package main

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

	"github.com/cesarFuhr/gocrypto/keys"
	"github.com/cesarFuhr/gocrypto/presenters"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwe"
)

type KeyStoreStub struct {
	CalledWith       []interface{}
	LastDeliveredKey keys.Key
}

var rsaKey, _ = rsa.GenerateKey(rand.Reader, 4098)

func (s *KeyStoreStub) CreateKey(scope string, exp time.Time) keys.Key {
	s.CalledWith = []interface{}{scope, exp}
	return keys.Key{
		Scope:      scope,
		Expiration: time.Now().AddDate(0, 0, 1),
		ID:         uuid.New().String(),
		Pub:        &rsaKey.PublicKey,
		Priv:       rsaKey,
	}
}

func (s *KeyStoreStub) FindKey(id string) (keys.Key, error) {
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

type CryptoStub struct {
	CalledWith []interface{}
}

func (c *CryptoStub) Encrypt(k *rsa.PublicKey, m string) ([]byte, error) {
	c.CalledWith = []interface{}{k, m}
	if m == "ERROR" {
		return []byte{}, errors.New("Any error at all")
	}
	msg, err := jwe.Encrypt([]byte(m), jwa.RSA_OAEP_256, k, jwa.A128CBC_HS256, jwa.NoCompress)
	if err != nil {
		fmt.Println("Error decrypting", err)
	}
	return msg, nil
}

var validReqBody, _ = json.Marshal(keyOpts{"scope", time.Now().UTC().Format(time.RFC3339)})

func TestPOSTKeys(t *testing.T) {
	keyStoreStub := KeyStoreStub{}
	cryptoStub := CryptoStub{}
	server := KeyServer{&keyStoreStub, &cryptoStub}
	t.Run("Should return 201 on /keys", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/keys", bytes.NewBuffer(validReqBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusCreated)
	})
	t.Run("Should return valid json on /keys", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/keys", bytes.NewBuffer(validReqBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		respBytes, _ := ioutil.ReadAll(response.Body)

		if !json.Valid(respBytes) {
			t.Errorf("got an invalid JSON %q", respBytes)
		}
	})
	t.Run("Should return all properties on /keys response", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/keys", bytes.NewBuffer(validReqBody))
		response := httptest.NewRecorder()

		wants := []string{"publicKey", "keyID", "expiration"}

		server.ServeHTTP(response, request)
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

		server.ServeHTTP(response, request)

		wantedExpiration, _ := time.Parse(time.RFC3339, expiration)
		assertInsideSlice(t, keyStoreStub.CalledWith, wantedExpiration)
		assertInsideSlice(t, keyStoreStub.CalledWith, scope)
	})
	t.Run("Should only return in the /keys endpoint", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/otherEndpoint", nil)
		response := httptest.NewRecorder()

		want := http.StatusNotFound

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, want)
	})
	t.Run("Should return a BadRequest if body is nil", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/keys", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

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

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusBadRequest)
		assertInsideJSON(t, response.Body, "message", "Invalid: expiration property format")
	})
}

func TestGETKeys(t *testing.T) {
	keyStoreStub := KeyStoreStub{}
	cryptoStub := CryptoStub{}
	server := KeyServer{&keyStoreStub, &cryptoStub}
	t.Run("Should return a 200 if it was a success", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/keys", nil)
		response := httptest.NewRecorder()

		want := http.StatusOK

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, want)
	})
	t.Run("Should return a Key if it was a success", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/keys", nil)
		response := httptest.NewRecorder()

		wants := []string{"publicKey", "keyID", "expiration"}

		server.ServeHTTP(response, request)
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
		server.ServeHTTP(response, request)
		respMap := map[string]interface{}{}
		extractJSON(response.Body, respMap)

		getRequest, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/keys?keyID=%v", respMap["keyID"]), nil)
		server.ServeHTTP(response, getRequest)

		assertInsideSlice(t, keyStoreStub.CalledWith, respMap["keyID"])
	})
	t.Run("If key was not found", func(t *testing.T) {
		t.Run("Should return a 404", func(t *testing.T) {
			want := http.StatusNotFound
			getRequest, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/keys?keyID=notFound"), nil)
			response := httptest.NewRecorder()
			server.ServeHTTP(response, getRequest)

			assertStatus(t, response.Code, want)
		})
		t.Run("Should a not found message", func(t *testing.T) {
			want := http.StatusNotFound
			getRequest, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/keys?keyID=notFound"), nil)
			response := httptest.NewRecorder()
			server.ServeHTTP(response, getRequest)

			assertStatus(t, response.Code, want)
			assertInsideJSON(t, response.Body, "message", "Key was not found")
		})
	})
	t.Run("If there was any other error", func(t *testing.T) {
		t.Run("Should return a 500", func(t *testing.T) {
			want := http.StatusInternalServerError
			getRequest, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/keys?keyID=otherError"), nil)
			response := httptest.NewRecorder()
			server.ServeHTTP(response, getRequest)

			assertStatus(t, response.Code, want)
		})
		t.Run("Should return a internal server error message", func(t *testing.T) {
			want := http.StatusInternalServerError
			getRequest, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/keys?keyID=otherError"), nil)
			response := httptest.NewRecorder()
			server.ServeHTTP(response, getRequest)

			assertStatus(t, response.Code, want)
			assertInsideJSON(t, response.Body, "message", "There was an unexpected error")
		})
	})
	t.Run("If uses a unsuported method", func(t *testing.T) {
		t.Run("Should return a method not allowed", func(t *testing.T) {
			want := http.StatusMethodNotAllowed
			getRequest, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/keys"), nil)
			response := httptest.NewRecorder()
			server.ServeHTTP(response, getRequest)

			assertStatus(t, response.Code, want)
			assertInsideJSON(t, response.Body, "message", "Method not allowed")
		})
	})
}

func TestEncrypt(t *testing.T) {
	keyStoreStub := KeyStoreStub{}
	cryptoStub := CryptoStub{}
	server := KeyServer{&keyStoreStub, &cryptoStub}
	t.Run("Should return a 200 if it was a success", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/encrypt", nil)
		response := httptest.NewRecorder()

		want := http.StatusOK

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, want)
	})
	t.Run("Should return the correct properties", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/encrypt", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var got presenters.HttpEncrypt
		json.NewDecoder(response.Body).Decode(&got)

		if got.EncryptedData == "" {
			t.Errorf("Expecting data prop, got %v", got.EncryptedData)
		}
	})
	t.Run("Should call findKeys with the right params", func(t *testing.T) {
		scope := "testing"
		expiration := time.Now().UTC().AddDate(0, 0, 1).Format(time.RFC3339)
		requestBody, _ := json.Marshal(map[string]string{
			"scope":      scope,
			"expiration": expiration,
		})

		request, _ := http.NewRequest(http.MethodPost, "/keys", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		respMap := map[string]interface{}{}
		extractJSON(response.Body, respMap)

		keyID := fmt.Sprintf("%v", respMap["keyID"])
		requestBody, _ = json.Marshal(map[string]string{
			"keyID": keyID,
		})
		request, _ = http.NewRequest(http.MethodPost, "/encrypt", bytes.NewBuffer(requestBody))
		server.ServeHTTP(response, request)

		assertInsideSlice(t, keyStoreStub.CalledWith, keyID)
	})
	t.Run("Should call Encrypt with the right params", func(t *testing.T) {
		scope := "testing"
		expiration := time.Now().UTC().AddDate(0, 0, 1).Format(time.RFC3339)
		requestBody, _ := json.Marshal(map[string]string{
			"scope":      scope,
			"expiration": expiration,
		})

		request, _ := http.NewRequest(http.MethodPost, "/keys", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		respMap := map[string]interface{}{}
		extractJSON(response.Body, respMap)

		keyID := fmt.Sprintf("%v", respMap["keyID"])
		data := "testing"
		requestBody, _ = json.Marshal(map[string]string{
			"keyID": keyID,
			"data":  data,
		})
		request, _ = http.NewRequest(http.MethodPost, "/encrypt", bytes.NewBuffer(requestBody))
		server.ServeHTTP(response, request)

		assertInsideSlice(t, cryptoStub.CalledWith, data)
		assertInsideSlice(t, cryptoStub.CalledWith, keyStoreStub.LastDeliveredKey.Pub)
	})
	t.Run("Should return a internal server error if there was a problem encrypting", func(t *testing.T) {
		scope := "testing"
		expiration := time.Now().UTC().AddDate(0, 0, 1).Format(time.RFC3339)
		requestBody, _ := json.Marshal(map[string]string{
			"scope":      scope,
			"expiration": expiration,
		})

		request, _ := http.NewRequest(http.MethodPost, "/keys", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		respMap := map[string]interface{}{}
		extractJSON(response.Body, respMap)

		keyID := fmt.Sprintf("%v", respMap["keyID"])
		data := "ERROR"
		requestBody, _ = json.Marshal(map[string]string{
			"keyID": keyID,
			"data":  data,
		})
		request, _ = http.NewRequest(http.MethodPost, "/encrypt", bytes.NewBuffer(requestBody))
		response = httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusInternalServerError)
		assertInsideJSON(t, response.Body, "message", "There was an unexpected error")
	})
	t.Run("Should return a valid jwe string", func(t *testing.T) {
		keyID := "id"
		requestBody, _ := json.Marshal(map[string]string{
			"keyID": keyID,
			"data":  "1234",
		})
		request, _ := http.NewRequest(http.MethodPost, "/encrypt", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		var e presenters.HttpEncrypt
		json.NewDecoder(response.Body).Decode(&e)

		got := e.EncryptedData

		if _, err := jwe.Decrypt([]byte(got), jwa.RSA_OAEP_256, rsaKey); err != nil {
			t.Errorf("response was not a valid jwe, %v", err)
		}
	})
	t.Run("Should return a precondition fail if the key does not exists", func(t *testing.T) {
		requestBody, _ := json.Marshal(map[string]string{
			"keyID": "notFound",
		})
		request, _ := http.NewRequest(http.MethodPost, "/encrypt", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusPreconditionFailed)
		assertInsideJSON(t, response.Body, "message", "Key was not found")
	})
	t.Run("Should return a internal server error if there is a error while finding key", func(t *testing.T) {
		requestBody, _ := json.Marshal(map[string]string{
			"keyID": "otherError",
		})
		request, _ := http.NewRequest(http.MethodPost, "/encrypt", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusInternalServerError)
		assertInsideJSON(t, response.Body, "message", "There was an unexpected error")
	})
	t.Run("If uses a unsuported method", func(t *testing.T) {
		t.Run("Should return a method not allowed", func(t *testing.T) {
			want := http.StatusMethodNotAllowed
			getRequest, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/encrypt"), nil)
			response := httptest.NewRecorder()
			server.ServeHTTP(response, getRequest)

			assertStatus(t, response.Code, want)
			assertInsideJSON(t, response.Body, "message", "Method not allowed")
		})
	})
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
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
