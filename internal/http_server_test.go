package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type keyHandlerStub struct {
	G struct {
		CalledWith []interface{}
	}
	P struct {
		CalledWith []interface{}
	}
}

func (h *keyHandlerStub) Post(w http.ResponseWriter, r *http.Request) {
	h.P.CalledWith = []interface{}{w, r}
}

func (h *keyHandlerStub) Get(w http.ResponseWriter, r *http.Request) {
	h.G.CalledWith = []interface{}{w, r}
}

func TestKeysEndpoint(t *testing.T) {
	kH := keyHandlerStub{}
	server := httpServer{&kH}
	t.Run("calls keyHandler.Post in a /keys http POST", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/keys", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertInsideSlice(t, kH.P.CalledWith, response)
		assertInsideSlice(t, kH.P.CalledWith, request)
	})
	t.Run("calls keyHandler.Get in a /keys http GET", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/keys", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertInsideSlice(t, kH.G.CalledWith, response)
		assertInsideSlice(t, kH.G.CalledWith, request)
	})
	t.Run("returns method not allowed for any other method", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPatch, "/keys", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertValue(t, response.Code, http.StatusMethodNotAllowed)
	})
}

func assertValue(t *testing.T, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Errorf("want %d, got %d", want, got)
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
