package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPOSTKeys(t *testing.T) {
	t.Run("Should return 200 on /keys", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/keys", nil)
		response := httptest.NewRecorder()

		KeyServer(response, request)

		assertStatus(t, response.Code, http.StatusCreated)
	})
	t.Run("Should return valid json on /keys", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/keys", nil)
		response := httptest.NewRecorder()

		KeyServer(response, request)

		respBytes, _ := ioutil.ReadAll(response.Body)

		if !json.Valid(respBytes) {
			t.Errorf("got an invalid JSON %q", respBytes)
		}
	})
	t.Run("Should have the publicKey property on /keys response", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/keys", nil)
		response := httptest.NewRecorder()

		KeyServer(response, request)
		respBytes, _ := ioutil.ReadAll(response.Body)
		respMap := map[string]interface{}{}
		_ = json.Unmarshal(respBytes, &respMap)

		if _, ok := respMap["publicKey"]; ok != true {
			t.Errorf("does not have the publicKey prop")
		}
	})
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
