package observability

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupTracing_NoEndpoint_IsNoOpAndDoesNotHang(t *testing.T) {
	shutdown, err := SetupTracing(context.Background(), TracingConfig{})
	require.NoError(t, err)
	require.NotNil(t, shutdown)

	// Starting a span against the no-op provider must be cheap and produce an
	// invalid (non-recording) span context.
	_, span := Tracer("test").Start(context.Background(), "noop")
	assert.False(t, span.SpanContext().IsValid(), "no-op provider must not record spans")
	span.End()

	// Shutdown returns immediately, never dialing a collector.
	assert.NoError(t, shutdown(context.Background()))
}

func TestSetupTracing_WithEndpoint_InstallsExporter(t *testing.T) {
	shutdown, err := SetupTracing(context.Background(), TracingConfig{
		Endpoint:    "localhost:4318",
		ServiceName: "test-svc",
	})
	require.NoError(t, err)
	require.NotNil(t, shutdown)

	// With an endpoint configured, the SDK provider samples, producing a valid
	// span context. We never actually export here.
	_, span := Tracer("test").Start(context.Background(), "real")
	assert.True(t, span.SpanContext().IsValid(), "configured provider should sample spans")
	span.End()

	assert.NoError(t, shutdown(context.Background()))
}

func TestNormalizeEndpoint(t *testing.T) {
	assert.Equal(t, "http://localhost:4318/v1/traces", normalizeEndpoint("localhost:4318"))
	assert.Equal(t, "https://collector.example.com", normalizeEndpoint("https://collector.example.com"))
	assert.Equal(t, "http://collector:4318/v1/traces", normalizeEndpoint("http://collector:4318/v1/traces"))
}

func TestMetrics_ObserveAndExpose(t *testing.T) {
	metrics := NewMetrics(nil)
	metrics.ObserveRequest("GET", "/teams/{id}", "200", 0.012)
	metrics.ObserveRequest("POST", "/auth/login", "401", 0.003)

	rec := httptest.NewRecorder()
	metrics.Handler().ServeHTTP(rec, httptest.NewRequest("GET", "/metrics", nil))

	body := rec.Body.String()
	assert.Equal(t, 200, rec.Code)
	assert.Contains(t, body, `http_requests_total{method="GET",route="/teams/{id}",status="200"} 1`)
	assert.Contains(t, body, `http_requests_total{method="POST",route="/auth/login",status="401"} 1`)
	assert.Contains(t, body, "http_request_duration_seconds_bucket")
	// No DB gauge when db is nil.
	assert.False(t, strings.Contains(body, "db_pool_in_use_connections"))
}
