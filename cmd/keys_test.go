package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cesarFuhr/gocrypto/internal/app/domain/keys"
	"github.com/google/uuid"
)

var (
	validCreateKeysBody, _ = json.Marshal(map[string]interface{}{
		"scope":      "test",
		"expiration": time.Now().UTC().Format(time.RFC3339),
	})
	invalidCreateKeysBody, _ = json.Marshal(map[string]interface{}{
		"scope":      "test",
		"expiration": time.Now().UTC().Format(time.Kitchen),
	})
	mockKey keys.Key
)

func TestKeys(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		endpoint     string
		body         []byte
		resultStatus int
	}{
		{
			"Creates and return a public key in der format",
			http.MethodPost,
			"/keys",
			validCreateKeysBody,
			201,
		},
		{
			"Returns a 400 bad request",
			http.MethodPost,
			"/keys",
			invalidCreateKeysBody,
			400,
		},
		{
			"Returns a 404 not found",
			http.MethodGet,
			"/keys/" + uuid.New().String(),
			nil,
			404,
		},
		{
			"Returns a public key in pem format",
			http.MethodGet,
			"/keys/" + mockKey.ID,
			nil,
			200,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var reqBody = bytes.NewReader([]byte{})

			if tt.body != nil {
				reqBody = bytes.NewReader(tt.body)
			}

			fmt.Println(tt.endpoint)
			request, _ := http.NewRequest(tt.method, tt.endpoint, reqBody)
			response := httptest.NewRecorder()

			httpServer.Handler.ServeHTTP(response, request)

			if response.Code != tt.resultStatus {
				t.Errorf("want %d, got %d", tt.resultStatus, response.Code)
			}
		})
	}
}
