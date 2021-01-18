package server

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

type keStub struct {
	G struct {
		CalledWith []interface{}
		Called     bool
	}
	P struct {
		CalledWith []interface{}
		Called     bool
	}
}

func (h *keStub) Post(w http.ResponseWriter, r *http.Request) {
	h.P.CalledWith = []interface{}{w, r}
	h.P.Called = true
}

func (h *keStub) Get(w http.ResponseWriter, r *http.Request) {
	h.G.CalledWith = []interface{}{w, r}
	h.G.Called = true
}

type encrypStub struct {
	P struct {
		CalledWith []interface{}
		Called     bool
	}
}

func (h *encrypStub) Post(w http.ResponseWriter, r *http.Request) {
	h.P.CalledWith = []interface{}{w, r}
	h.P.Called = true
}

type decrypStub struct {
	P struct {
		CalledWith []interface{}
		Called     bool
	}
}

func (h *decrypStub) Post(w http.ResponseWriter, r *http.Request) {
	h.P.CalledWith = []interface{}{w, r}
	h.P.Called = true
}

type loggerStub struct {
	CalledWith []interface{}
	Called     bool
}

func (l *loggerStub) Info(msg string, args ...zap.Field) {
	l.CalledWith = append(l.CalledWith, msg)
	l.CalledWith = append(l.CalledWith, args)
	l.Called = true
}

var (
	log    = new(loggerStub)
	kH     = new(keStub)
	eH     = new(encrypStub)
	dH     = new(decrypStub)
	server = NewHTTPServer(log, kH, eH, dH).Handler
)

func TestKeysEndpoint(t *testing.T) {
	t.Run("calls key.Post in a /keys http POST", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/keys", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertValue(t, kH.P.Called, true)
		kH.P.Called = false
	})
	t.Run("calls key.Get in a /keys http GET", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/keys/100", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertValue(t, kH.G.Called, true)
		kH.G.Called = false
	})
	t.Run("returns method not allowed for any other method", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPatch, "/keys", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertValue(t, response.Code, http.StatusMethodNotAllowed)
	})
	t.Run("calls logger.Info in /keys http Requests", func(t *testing.T) {
		log.CalledWith = []interface{}{}
		endpoint := "/keys"
		method := http.MethodPost
		request, _ := http.NewRequest(method, endpoint, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertValue(t, log.Called, true)
		log.Called = false
	})
}

func TestEncryptionEndpoint(t *testing.T) {
	t.Run("calls encryp.Post in a /encrypt http POST", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/encrypt", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertValue(t, eH.P.Called, true)
		eH.P.Called = false
	})
	t.Run("returns method not allowed for any other method", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPatch, "/encrypt", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertValue(t, response.Code, http.StatusMethodNotAllowed)
	})
	t.Run("calls logger.Info in /encrypt http Requests", func(t *testing.T) {
		log.CalledWith = []interface{}{}
		endpoint := "/encrypt"
		method := http.MethodPost
		request, _ := http.NewRequest(method, endpoint, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertValue(t, log.Called, true)
		log.Called = false
	})
}

func TestDecryptionEndpoint(t *testing.T) {
	t.Run("calls decryp.Post in a /decrypt http POST", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/decrypt", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertValue(t, dH.P.Called, true)
		dH.P.Called = false
	})
	t.Run("returns method not allowed for any other method", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPatch, "/decrypt", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertValue(t, response.Code, http.StatusMethodNotAllowed)
	})
	t.Run("calls logger.Info in /decrypt http Requests", func(t *testing.T) {
		log.CalledWith = []interface{}{}
		endpoint := "/decrypt"
		method := http.MethodPost
		request, _ := http.NewRequest(method, endpoint, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertValue(t, log.Called, true)
		log.Called = false
	})
}

func assertValue(t *testing.T, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Errorf("want %d, got %d", want, got)
	}
}

func assertType(t *testing.T, got, want interface{}) {
	t.Helper()
	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf("want %T, got %T", want, got)
	}
}
