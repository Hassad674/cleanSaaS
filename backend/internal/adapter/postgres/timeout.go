package postgres

import (
	"context"
	"time"

	"github.com/hassad/boilerplateSaaS/backend/pkg/ctxutil"
)

// defaultDBTimeout is the ceiling applied by ctxWithTimeout to a database
// operation that arrives without a nearer deadline. It guards background/no-
// deadline callers (scheduler cleanups, internal maintenance) from a query that
// hangs indefinitely. Request-scoped callers already carry the HTTP request
// deadline, which is nearer and therefore wins (see ctxutil.WithTimeout).
const defaultDBTimeout = 15 * time.Second

// ctxWithTimeout derives a context bounded by defaultDBTimeout unless the caller
// already set a nearer deadline. It is the single, shared place repositories use
// to bound a query, rather than scattering raw context.WithTimeout calls.
//
// Light-touch by design: it is NOT mandatory on every repository method (forcing
// it everywhere risks clamping legitimate long request flows and bloats the
// diff). Instead it is the building block for any path that needs a guaranteed
// ceiling. Today the two unbounded paths called out in the task — the scheduler
// cleanups and webhook processing — are already bounded upstream (the scheduler
// gives each job a deadline; the webhook runs under the HTTP request context),
// so applying it again at the repo layer would be redundant. The helper exists
// so future background callers have a clean, consistent way to opt in.
//
// Callers must always call the returned CancelFunc (defer cancel()).
func ctxWithTimeout(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	if d <= 0 {
		d = defaultDBTimeout
	}
	return ctxutil.WithTimeout(ctx, d)
}
