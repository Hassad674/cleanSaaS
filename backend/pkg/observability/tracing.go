// Package observability wires application telemetry — distributed tracing via
// OpenTelemetry — without leaking the SDK into the rest of the codebase.
//
// It is a leaf utility under pkg/: it imports only the Go standard library and
// the OpenTelemetry SDK, so it can be used by the composition root (cmd/api)
// and the HTTP layer (handler/middleware) without violating the hexagonal
// dependency rule.
//
// Design goal: zero friction for local runs. When no OTLP endpoint is
// configured the tracer provider is a never-sampling, no-op provider — it
// records nothing, exports nothing, never dials a collector, and never blocks
// on shutdown. Tracing is opt-in by setting OTEL_EXPORTER_OTLP_ENDPOINT.
package observability

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

// defaultServiceName is used when the caller passes an empty service name.
const defaultServiceName = "cleansaas-backend"

// TracingConfig configures the tracer provider. An empty Endpoint disables
// exporting entirely (no-op provider).
type TracingConfig struct {
	// Endpoint is the OTLP/HTTP collector endpoint (e.g. "localhost:4318" or a
	// full URL). Empty means tracing is disabled and a no-op provider is used.
	Endpoint string
	// ServiceName labels emitted spans. Defaults to "cleansaas-backend".
	ServiceName string
}

// Shutdown gracefully flushes and stops the tracer provider. It is always
// safe to call (the no-op case returns nil immediately).
type Shutdown func(ctx context.Context) error

// SetupTracing installs a global OpenTelemetry tracer provider and propagator.
//
// When cfg.Endpoint is empty it installs a no-op provider: nothing is sampled,
// nothing is exported, and the returned Shutdown is a cheap no-op. This keeps
// local development collector-free.
//
// When cfg.Endpoint is set it installs an SDK provider with a batching OTLP/HTTP
// exporter. The returned Shutdown flushes pending spans and tears the provider
// down; the caller should invoke it during graceful shutdown.
//
// The W3C trace-context + baggage propagator is always installed so inbound
// trace headers are honored regardless of whether we export.
func SetupTracing(ctx context.Context, cfg TracingConfig) (Shutdown, error) {
	// Always honor inbound W3C trace context, even in no-op mode, so request
	// IDs / trace IDs minted upstream are preserved through this service.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	if cfg.Endpoint == "" {
		// No collector configured: install an explicit no-op provider so callers
		// can start spans freely with zero cost and zero network activity.
		otel.SetTracerProvider(noop.NewTracerProvider())
		return func(context.Context) error { return nil }, nil
	}

	name := cfg.ServiceName
	if name == "" {
		name = defaultServiceName
	}

	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(normalizeEndpoint(cfg.Endpoint)))
	if err != nil {
		return nil, err
	}

	// NewSchemaless avoids a schema-URL conflict with resource.Default() (whose
	// schema version tracks the SDK, not our pinned semconv import).
	res, err := resource.Merge(
		resource.Default(),
		resource.NewSchemaless(semconv.ServiceName(name)),
	)
	if err != nil {
		return nil, err
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter, sdktrace.WithBatchTimeout(5*time.Second)),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample())),
	)

	otel.SetTracerProvider(provider)

	return provider.Shutdown, nil
}

// normalizeEndpoint accepts either a bare host:port or a full URL and returns a
// full URL suitable for otlptracehttp.WithEndpointURL. A bare host is assumed to
// speak plaintext HTTP on the standard OTLP traces path.
func normalizeEndpoint(endpoint string) string {
	if hasScheme(endpoint) {
		return endpoint
	}
	return "http://" + endpoint + "/v1/traces"
}

// hasScheme reports whether s begins with "http://" or "https://".
func hasScheme(s string) bool {
	const httpScheme = "http://"
	const httpsScheme = "https://"
	if len(s) >= len(httpScheme) && s[:len(httpScheme)] == httpScheme {
		return true
	}
	if len(s) >= len(httpsScheme) && s[:len(httpsScheme)] == httpsScheme {
		return true
	}
	return false
}

// Tracer returns a named tracer from the globally-installed provider. Callers in
// the HTTP layer use this to start server spans.
func Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}
