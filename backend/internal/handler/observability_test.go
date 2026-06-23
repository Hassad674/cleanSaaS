package handler

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hassad/boilerplateSaaS/backend/pkg/observability"
)

// --- fake sql driver -------------------------------------------------------
//
// A minimal database/sql driver whose Ping result is controllable, so /readyz
// can be exercised without a real database. database/sql calls Conn.Ping when
// the conn implements driver.Pinger, which our pingConn does.

var errPingFailed = errors.New("connection refused")

type pingDriver struct{ ok *bool }

func (d pingDriver) Open(string) (driver.Conn, error) { return pingConn{ok: d.ok}, nil }

type pingConn struct{ ok *bool }

func (c pingConn) Ping(context.Context) error {
	if c.ok != nil && *c.ok {
		return nil
	}
	return errPingFailed
}

func (c pingConn) Prepare(string) (driver.Stmt, error) { return nil, errors.New("not supported") }
func (c pingConn) Close() error                        { return nil }
func (c pingConn) Begin() (driver.Tx, error)           { return nil, errors.New("not supported") }

var (
	registerOnce sync.Once
	pingState    = new(bool)
)

// newPingDB returns a *sql.DB whose PingContext succeeds when healthy is true.
func newPingDB(t *testing.T, healthy bool) *sql.DB {
	t.Helper()
	registerOnce.Do(func() {
		sql.Register("pingfake", pingDriver{ok: pingState})
	})
	*pingState = healthy
	db, err := sql.Open("pingfake", "")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	return db
}

// newTestRouter builds the real router with only the dependencies the
// observability/health endpoints touch; optional services are nil (their routes
// are simply not registered).
func newTestRouter(t *testing.T, db *sql.DB, metrics *observability.Metrics) http.Handler {
	t.Helper()
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	resolver := func(_ context.Context, userID, _ string) (string, error) { return userID, nil }
	return NewRouter(
		nil, nil, nil, nil, nil, nil, nil, nil, // app services (unused by these routes)
		nil,         // wsHub
		"secret",    // jwtSecret
		"http://fe", // frontendURL
		db, logger,
		nil,      // demoAI
		resolver, // orgResolver
		metrics,
	)
}

func TestLivez_AlwaysOK_NoDependencyCheck(t *testing.T) {
	// DB is unhealthy, but liveness must NOT check dependencies.
	router := newTestRouter(t, newPingDB(t, false), nil)

	req := httptest.NewRequest(http.MethodGet, "/livez", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "alive")
}

func TestReadyz_OK_WhenDBReachable(t *testing.T) {
	router := newTestRouter(t, newPingDB(t, true), nil)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"status":"ready"`)
}

func TestReadyz_503_WhenDBPingFails(t *testing.T) {
	router := newTestRouter(t, newPingDB(t, false), nil)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.Contains(t, rec.Body.String(), "not ready")
}

func TestHealth_BackwardsCompatPayload(t *testing.T) {
	router := newTestRouter(t, newPingDB(t, true), nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, `"status":"ok"`)
	assert.Contains(t, body, `"version"`)
	assert.Contains(t, body, `"uptime"`)
}

func TestMetrics_ExposesSeriesAfterRequest(t *testing.T) {
	metrics := observability.NewMetrics(newPingDB(t, true))
	router := newTestRouter(t, newPingDB(t, true), metrics)

	// Drive one request through a real (matched) route so a series is recorded.
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/readyz", nil))
	require.Equal(t, http.StatusOK, rec.Code)

	// Scrape /metrics.
	scrape := httptest.NewRecorder()
	router.ServeHTTP(scrape, httptest.NewRequest(http.MethodGet, "/metrics", nil))

	assert.Equal(t, http.StatusOK, scrape.Code)
	body := scrape.Body.String()
	assert.Contains(t, body, "http_requests_total")
	assert.Contains(t, body, "http_request_duration_seconds")
	assert.Contains(t, body, "db_pool_in_use_connections")
	// The recorded request used the route PATTERN as the label, not the raw path.
	assert.Contains(t, body, `route="/readyz"`)
}

func TestMetrics_DoesNotCountItself(t *testing.T) {
	metrics := observability.NewMetrics(nil)
	router := newTestRouter(t, newPingDB(t, true), metrics)

	// Scrape twice; /metrics must not appear as its own series.
	for i := 0; i < 2; i++ {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))
		require.Equal(t, http.StatusOK, rec.Code)
	}

	scrape := httptest.NewRecorder()
	router.ServeHTTP(scrape, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	assert.NotContains(t, scrape.Body.String(), `route="/metrics"`)
}

func TestRequestID_Generated_WhenAbsent(t *testing.T) {
	router := newTestRouter(t, newPingDB(t, true), nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/livez", nil))

	assert.NotEmpty(t, rec.Header().Get("X-Request-ID"), "a request id should be generated and echoed")
}

func TestRequestID_Preserved_WhenInbound(t *testing.T) {
	router := newTestRouter(t, newPingDB(t, true), nil)

	const inbound = "client-supplied-123"
	req := httptest.NewRequest(http.MethodGet, "/livez", nil)
	req.Header.Set("X-Request-ID", inbound)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, inbound, rec.Header().Get("X-Request-ID"), "inbound request id must be preserved")
}
