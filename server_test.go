package main

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/cesarFuhr/gocrypto/keys"
)

type KeyStoreStub struct {
	CalledWith []interface{}
}

func (s *KeyStoreStub) CreateKey(scope string, exp time.Time) keys.Key {
	s.CalledWith = []interface{}{scope, exp}
	return keys.Key{
		Scope:      scope,
		Expiration: time.Now().AddDate(0, 0, 1),
		ID:         "a key ID",
		Pub:        &rsa.PublicKey{},
		Priv:       &rsa.PrivateKey{},
	}
}

var validReqBody, _ = json.Marshal(keyOpts{"scope", time.Now().UTC().Format(time.RFC3339)})

func TestPOSTKeys(t *testing.T) {
	keyStoreStub := KeyStoreStub{}
	server := KeyServer{&keyStoreStub}
	t.Run("Should return 200 on /keys", func(t *testing.T) {
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
		extractJson(response.Body, respMap)

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
			"scope": scope,
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
	t.Run("Should reutrn a BadRequest if body is nil", func(t *testing.T){
		request, _ := http.NewRequest(http.MethodPost, "/keys", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusBadRequest)
		assertInsideJson(t, response.Body, "message", "Invalid: Empty body")
	})
	t.Run("Should reutrn a BadRequest if expiration is not a RFC3339", func(t *testing.T){
		scope := "testing"
		expiration := time.Now().UTC().AddDate(0, 0, 1).Format(time.RubyDate)
		requestBody, _ := json.Marshal(map[string]string{
			"scope": scope,
			"expiration": expiration,
		})
		request, _ := http.NewRequest(http.MethodPost, "/keys", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusBadRequest)
		assertInsideJson(t, response.Body, "message", "Invalid: expiration property format")
	})
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func extractJson(jBuff *bytes.Buffer, m map[string]interface{}) error {
	respBytes, _ := ioutil.ReadAll(jBuff)
	_ = json.Unmarshal(respBytes, &m)
	return nil
}

func assertInsideJson(t *testing.T, jBuff *bytes.Buffer, wantedKey string, wantedValue interface{}) {
	t.Helper()
	got := map[string]interface{}{}
	extractJson(jBuff, got)
	if !reflect.DeepEqual(got[wantedKey], wantedValue) {
		t.Errorf("got %v, want %v", got[wantedKey], wantedValue)
	}
}

func assertInsideSlice(t *testing.T, a []interface{}, want interface{}) {
    t.Helper()
	has := false
	for _, v := range a {
		if reflect.DeepEqual(v, want) {
			has = true
		}
	}
	if !has {
		t.Errorf("Did not found: %v, in %v", want, a)
	}
}
