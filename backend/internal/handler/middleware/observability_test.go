package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hassad/boilerplateSaaS/backend/pkg/observability"
)

func TestStructuredLogging_LevelByStatus(t *testing.T) {
	cases := []struct {
		status    int
		wantLevel string
	}{
		{http.StatusOK, "INFO"},
		{http.StatusBadRequest, "WARN"},
		{http.StatusForbidden, "WARN"},
		{http.StatusInternalServerError, "ERROR"},
	}

	for _, tc := range cases {
		var buf bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&buf, nil))
		handler := StructuredLogging(logger)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(tc.status)
		}))

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))

		var line map[string]any
		require.NoError(t, json.Unmarshal(buf.Bytes(), &line))
		assert.Equal(t, tc.wantLevel, line["level"], "status %d", tc.status)
	}
}

func TestStructuredLogging_IncludesUserAndOrg(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	// Logging installs the holder; an inner middleware (mimicking Auth) fills it.
	inner := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	withUser := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			setLogUser(r.Context(), "user-1", "org-9")
			next.ServeHTTP(w, r)
		})
	}
	handler := StructuredLogging(logger)(withUser(inner))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))

	var line map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &line))
	assert.Equal(t, "user-1", line["user_id"])
	assert.Equal(t, "org-9", line["org"])
}

func TestMetricsMiddleware_UsesRoutePattern(t *testing.T) {
	metrics := observability.NewMetrics(nil)

	r := chi.NewRouter()
	r.Use(Metrics(metrics))
	r.Get("/teams/{id}", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/teams/abc-123", nil))
	require.Equal(t, http.StatusOK, rec.Code)

	scrape := httptest.NewRecorder()
	metrics.Handler().ServeHTTP(scrape, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	body := scrape.Body.String()

	// The label is the pattern, never the raw path (high-cardinality guard).
	assert.Contains(t, body, `route="/teams/{id}"`)
	assert.NotContains(t, body, "abc-123")
}

func TestRequestID_RoundTrip(t *testing.T) {
	var captured string
	h := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = RequestIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	// Generated when absent.
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))
	assert.NotEmpty(t, captured)
	assert.Equal(t, captured, rec.Header().Get("X-Request-ID"))

	// Preserved when inbound.
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("X-Request-ID", "inbound-abc")
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, req)
	assert.Equal(t, "inbound-abc", captured)
	assert.Equal(t, "inbound-abc", rec2.Header().Get("X-Request-ID"))
}
