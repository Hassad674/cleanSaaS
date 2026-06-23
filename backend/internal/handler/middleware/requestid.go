package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const requestIDKey contextKey = "requestID"

// requestIDHeader is the canonical header used to carry a request ID across the
// trust boundary (e.g. from an upstream gateway or the client).
const requestIDHeader = "X-Request-ID"

// RequestIDFromContext returns the request ID stored in ctx, or "" if none.
func RequestIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(requestIDKey).(string)
	return v
}

// RequestID is middleware that establishes a per-request correlation ID.
//
// It honors an inbound X-Request-ID header when present (so a trace started
// upstream is preserved); otherwise it mints a new UUID. The ID is stored in the
// request context (read via RequestIDFromContext) AND echoed back on the
// response as X-Request-ID so clients and logs can correlate.
//
// It runs before StructuredLogging so the logging middleware can read the ID
// from context.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(requestIDHeader)
		if reqID == "" {
			reqID = uuid.New().String()
		}

		// Echo it back immediately so it is present even on early returns
		// (e.g. a rate-limit rejection downstream).
		w.Header().Set(requestIDHeader, reqID)

		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
