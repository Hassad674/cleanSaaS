package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCORS_SetsHeaders_MatchingOrigin(t *testing.T) {
	handler := CORS("http://localhost:3006")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set("Origin", "http://localhost:3006")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, "http://localhost:3006", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Methods"), "GET")
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Methods"), "POST")
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Methods"), "PUT")
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Methods"), "DELETE")
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Headers"), "Authorization")
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
	assert.Equal(t, "86400", rec.Header().Get("Access-Control-Max-Age"))
	assert.Equal(t, "true", rec.Header().Get("Access-Control-Allow-Credentials"),
		"credentials should be set for matching origin")
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCORS_NonMatchingOrigin_NoAllowOrigin(t *testing.T) {
	handler := CORS("http://localhost:3006")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set("Origin", "http://evil.com")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Origin"),
		"non-matching origin should not get Allow-Origin header")
	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Credentials"),
		"non-matching origin should not get credentials header")
	// Methods and headers are still set (general headers)
	assert.NotEmpty(t, rec.Header().Get("Access-Control-Allow-Methods"))
}

func TestCORS_NoOriginHeader(t *testing.T) {
	handler := CORS("http://localhost:3006")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	// No Origin header set
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Origin"),
		"missing Origin header should result in no Allow-Origin")
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCORS_OptionsPreflightReturns204(t *testing.T) {
	var handlerCalled bool
	handler := CORS("http://localhost:3006")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	}))

	req := httptest.NewRequest(http.MethodOptions, "/api/test", nil)
	req.Header.Set("Origin", "http://localhost:3006")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code, "OPTIONS should return 204")
	assert.False(t, handlerCalled, "next handler should not be called for OPTIONS")
	assert.Equal(t, "http://localhost:3006", rec.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_MultipleOrigins(t *testing.T) {
	handler := CORS("http://localhost:3006", "https://myapp.com")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// First origin
	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set("Origin", "http://localhost:3006")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, "http://localhost:3006", rec.Header().Get("Access-Control-Allow-Origin"))

	// Second origin
	req2 := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req2.Header.Set("Origin", "https://myapp.com")
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)
	assert.Equal(t, "https://myapp.com", rec2.Header().Get("Access-Control-Allow-Origin"))

	// Unknown origin
	req3 := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req3.Header.Set("Origin", "https://unknown.com")
	rec3 := httptest.NewRecorder()
	handler.ServeHTTP(rec3, req3)
	assert.Empty(t, rec3.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_PassesThroughNonOptionsRequests(t *testing.T) {
	var handlerCalled bool
	handler := CORS("http://localhost:3006")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusCreated)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/test", nil)
	req.Header.Set("Origin", "http://localhost:3006")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.True(t, handlerCalled, "non-OPTIONS requests should pass through to handler")
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestCORS_VaryHeaderSet(t *testing.T) {
	handler := CORS("http://localhost:3006")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set("Origin", "http://localhost:3006")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, "Origin", rec.Header().Get("Vary"),
		"Vary header should be set to Origin for proper caching")
}

func TestCORS_TrailingSlashNormalization(t *testing.T) {
	handler := CORS("http://localhost:3006/")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set("Origin", "http://localhost:3006")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, "http://localhost:3006", rec.Header().Get("Access-Control-Allow-Origin"),
		"trailing slash should be normalized to match")
}
