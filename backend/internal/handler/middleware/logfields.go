package middleware

import "context"

// logFields is a mutable, per-request holder for correlation attributes that are
// only known AFTER an inner middleware runs (e.g. user_id and org are resolved
// by Auth, which sits below the logging middleware in the chain).
//
// Because context values propagate only downward, the outermost logging
// middleware cannot observe a context that inner layers replaced. Instead it
// installs a pointer to this struct up front; inner layers mutate the pointee,
// and the logger reads it when the request completes. Access is single-threaded
// per request (no concurrent goroutines touch it), so no locking is needed.
type logFields struct {
	userID string
	org    string
}

const logFieldsKey contextKey = "logFields"

// withLogFields returns a context carrying a fresh, empty logFields holder.
func withLogFields(ctx context.Context) (context.Context, *logFields) {
	lf := &logFields{}
	return context.WithValue(ctx, logFieldsKey, lf), lf
}

// logFieldsFromContext returns the holder installed by the logging middleware,
// or nil if logging is not in the chain (e.g. in a unit test exercising a
// handler in isolation).
func logFieldsFromContext(ctx context.Context) *logFields {
	lf, _ := ctx.Value(logFieldsKey).(*logFields)
	return lf
}

// setLogUser records the authenticated user and org on the request's log holder,
// if one is present. Called by Auth so 4xx/5xx logs can be attributed.
func setLogUser(ctx context.Context, userID, org string) {
	if lf := logFieldsFromContext(ctx); lf != nil {
		lf.userID = userID
		lf.org = org
	}
}
