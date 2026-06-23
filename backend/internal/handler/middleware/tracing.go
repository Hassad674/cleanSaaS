package middleware

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

// tracerName identifies the instrumentation scope for HTTP server spans.
const tracerName = "github.com/hassad/boilerplateSaaS/backend/internal/handler"

// Tracing returns middleware that starts one server span per request and
// propagates inbound W3C trace context, so this service joins any distributed
// trace started upstream.
//
// The span is named by the Chi route PATTERN (resolved after routing) to keep
// span names low-cardinality. When no OTLP endpoint is configured the global
// provider is a no-op, so this middleware is effectively free and never dials a
// collector.
func Tracing(next http.Handler) http.Handler {
	tracer := otel.Tracer(tracerName)
	propagator := otel.GetTextMapPropagator()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		ctx, span := tracer.Start(ctx, r.Method,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.HTTPRequestMethodKey.String(r.Method),
				semconv.URLPath(r.URL.Path),
			),
		)
		defer span.End()

		wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(wrapped, r.WithContext(ctx))

		// Refine the span name now that routing has matched a pattern; this keeps
		// the span name space bounded (one name per route, not per path).
		pattern := routePattern(r)
		span.SetName(r.Method + " " + pattern)
		span.SetAttributes(
			attribute.String("http.route", pattern),
			semconv.HTTPResponseStatusCode(wrapped.status),
		)
		if wrapped.status >= http.StatusInternalServerError {
			span.SetStatus(codes.Error, http.StatusText(wrapped.status))
		}
	})
}
