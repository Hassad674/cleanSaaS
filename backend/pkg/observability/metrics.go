package observability

import (
	"database/sql"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics owns the application's Prometheus collectors on a DEDICATED registry
// (not the global default), so /metrics exposes exactly our series and nothing a
// transitively-imported library might have registered globally.
type Metrics struct {
	registry        *prometheus.Registry
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	dbInUse         prometheus.GaugeFunc
}

// NewMetrics builds the collectors and registers them on a fresh registry.
//
// If db is non-nil, a gauge reporting the number of in-use pooled connections is
// registered. It is a GaugeFunc, so it reads sql.DB.Stats() lazily at scrape
// time — no background goroutine, no polling cost between scrapes.
func NewMetrics(db *sql.DB) *Metrics {
	registry := prometheus.NewRegistry()

	requestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed, labeled by method, route pattern and status code.",
		},
		[]string{"method", "route", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds, labeled by method and route pattern.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route"},
	)

	registry.MustRegister(requestsTotal, requestDuration)

	metrics := &Metrics{
		registry:        registry,
		requestsTotal:   requestsTotal,
		requestDuration: requestDuration,
	}

	if db != nil {
		metrics.dbInUse = prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name: "db_pool_in_use_connections",
				Help: "Number of database connections currently in use (sql.DB.Stats().InUse).",
			},
			func() float64 { return float64(db.Stats().InUse) },
		)
		registry.MustRegister(metrics.dbInUse)
	}

	return metrics
}

// ObserveRequest records one finished HTTP request. route should be the low
// cardinality route PATTERN (e.g. "/teams/{id}"), never the raw path, to keep
// the series count bounded.
func (m *Metrics) ObserveRequest(method, route, status string, seconds float64) {
	m.requestsTotal.WithLabelValues(method, route, status).Inc()
	m.requestDuration.WithLabelValues(method, route).Observe(seconds)
}

// Handler returns the HTTP handler that serves the Prometheus exposition format
// for THIS registry only. Mount it (without auth) at /metrics.
func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}
