package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/hassad/boilerplateSaaS/backend/pkg/observability"
)

// metricsExcludedPaths are endpoints we never count, to keep the scrape target
// from observing itself (and the cheap liveness probe out of the latency stats).
var metricsExcludedPaths = map[string]bool{
	"/metrics": true,
	"/livez":   true,
}

// Metrics returns middleware that records HTTP request count and latency into
// the provided collector set, labeled by the low-cardinality Chi route PATTERN
// (e.g. "/teams/{id}") rather than the raw path.
//
// /metrics and /livez are excluded so the scrape endpoint and the high-frequency
// liveness probe do not pollute the series. It must sit below routing-agnostic
// middleware but is added globally; the route pattern is resolved after the
// inner handler matches a route.
func Metrics(collector *observability.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if metricsExcludedPaths[r.URL.Path] {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			collector.ObserveRequest(
				r.Method,
				routePattern(r),
				strconv.Itoa(wrapped.status),
				time.Since(start).Seconds(),
			)
		})
	}
}

// routePattern returns the matched Chi route pattern for r, falling back to
// "unmatched" when no route matched (404) so unknown paths cannot explode the
// label cardinality.
func routePattern(r *http.Request) string {
	if rc := chi.RouteContext(r.Context()); rc != nil {
		if pattern := rc.RoutePattern(); pattern != "" {
			return pattern
		}
	}
	return "unmatched"
}
