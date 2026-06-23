package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructuredLogging_LogsRequest(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	handler := StructuredLogging(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	logged := buf.String()
	assert.Contains(t, logged, "request_id")
	assert.Contains(t, logged, "/health")
	assert.Contains(t, logged, "GET")
}

func TestStructuredLogging_AddsRequestID(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	// RequestID mints/echoes the correlation ID; StructuredLogging reads it from
	// the context. They are chained as they are in the real router.
	var capturedID string
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = RequestIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})
	handler := RequestID(StructuredLogging(logger)(inner))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.NotEmpty(t, capturedID, "request ID should be set in context")
	assert.Equal(t, capturedID, rec.Header().Get("X-Request-ID"), "request ID should be echoed on the response")
}
