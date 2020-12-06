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

type encryptHandlerStub struct {
	P struct {
		CalledWith []interface{}
	}
}

func (h *encryptHandlerStub) Post(w http.ResponseWriter, r *http.Request) {
	h.P.CalledWith = []interface{}{w, r}
}

type decryptHandlerStub struct {
	P struct {
		CalledWith []interface{}
	}
}

func (h *decryptHandlerStub) Post(w http.ResponseWriter, r *http.Request) {
	h.P.CalledWith = []interface{}{w, r}
}

type loggerStub struct {
	CalledWith []interface{}
}

func (l *loggerStub) Info(args ...interface{}) {
	l.CalledWith = args
}

var (
	log    = loggerStub{}
	kH     = keyHandlerStub{}
	eH     = encryptHandlerStub{}
	dH     = decryptHandlerStub{}
	server = httpServer{&log, &kH, &eH, &dH}
)

func TestKeysEndpoint(t *testing.T) {
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
	t.Run("calls logger.Info in /keys http Requests", func(t *testing.T) {
		log.CalledWith = []interface{}{}
		endpoint := "/keys"
		method := http.MethodPost
		request, _ := http.NewRequest(method, endpoint, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertInsideSlice(t, log.CalledWith, endpoint)
		assertInsideSlice(t, log.CalledWith, method)
	})
}

func TestEncryptionEndpoint(t *testing.T) {
	t.Run("calls encryptHandler.Post in a /encrypt http POST", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/encrypt", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertInsideSlice(t, eH.P.CalledWith, response)
		assertInsideSlice(t, eH.P.CalledWith, request)
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

		assertInsideSlice(t, log.CalledWith, endpoint)
		assertInsideSlice(t, log.CalledWith, method)
	})
}

func TestDecryptionEndpoint(t *testing.T) {
	t.Run("calls decryptHandler.Post in a /decrypt http POST", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/decrypt", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertInsideSlice(t, dH.P.CalledWith, response)
		assertInsideSlice(t, dH.P.CalledWith, request)
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

		assertInsideSlice(t, log.CalledWith, endpoint)
		assertInsideSlice(t, log.CalledWith, method)
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
