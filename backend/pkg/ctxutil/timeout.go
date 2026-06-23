// Package ctxutil provides small, dependency-free context helpers shared across
// adapters. It imports only the Go standard library so it can be used from any
// layer without violating the hexagonal dependency rule.
package ctxutil

import (
	"context"
	"time"
)

// WithTimeout returns a derived context bounded by d, unless the incoming ctx
// already carries a nearer deadline (in which case the existing deadline wins
// and the returned cancel func is a no-op cleanup of the new context).
//
// This lets adapter boundaries impose a sane default ceiling on a hung external
// call without ever EXTENDING a deadline the caller already set. A non-positive
// d means "no timeout": the original ctx is returned unchanged.
//
// Callers must always call the returned CancelFunc (defer cancel()) to release
// resources, exactly as with context.WithTimeout.
func WithTimeout(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	if d <= 0 {
		return ctx, func() {}
	}

	// If the caller's context already expires sooner than our proposed default,
	// respect it — do not push the deadline out.
	if existing, ok := ctx.Deadline(); ok {
		if time.Until(existing) <= d {
			return ctx, func() {}
		}
	}

	return context.WithTimeout(ctx, d)
}
