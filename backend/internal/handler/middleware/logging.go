package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	status int
	wrote  bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wrote {
		rw.status = code
		rw.wrote = true
	}
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wrote {
		rw.status = http.StatusOK
		rw.wrote = true
	}
	return rw.ResponseWriter.Write(b)
}

// Flush implements http.Flusher so SSE streaming works through this wrapper.
func (rw *responseWriter) Flush() {
	if f, ok := rw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// StructuredLogging logs each request with slog in JSON format.
//
// Severity reflects the outcome: 5xx logs at Error, 4xx at Warn, everything
// else at Info. Every line carries the request_id (set by RequestID, which must
// run before this), and — when an authenticated route resolved them — user_id
// and org. It is stateless and cheap: a single log call per request.
func StructuredLogging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ctx, fields := withLogFields(r.Context())
			r = r.WithContext(ctx)

			wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(wrapped, r)

			logRequest(logger, r, wrapped.status, time.Since(start), fields)
		})
	}
}

// logRequest emits the structured access-log line at a severity matching status.
func logRequest(logger *slog.Logger, r *http.Request, status int, duration time.Duration, fields *logFields) {
	attrs := []slog.Attr{
		slog.String("request_id", RequestIDFromContext(r.Context())),
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.Int("status", status),
		slog.Duration("duration", duration),
		slog.String("remote_addr", r.RemoteAddr),
		slog.String("user_agent", r.UserAgent()),
	}
	if fields != nil {
		if fields.userID != "" {
			attrs = append(attrs, slog.String("user_id", fields.userID))
		}
		if fields.org != "" {
			attrs = append(attrs, slog.String("org", fields.org))
		}
	}

	logger.LogAttrs(r.Context(), levelForStatus(status), "request", attrs...)
}

// levelForStatus maps an HTTP status code to a log severity: 5xx → Error,
// 4xx → Warn, otherwise Info.
func levelForStatus(status int) slog.Level {
	switch {
	case status >= http.StatusInternalServerError:
		return slog.LevelError
	case status >= http.StatusBadRequest:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}
