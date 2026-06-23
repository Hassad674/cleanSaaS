// Package orgctx carries the caller's active organization ID through a
// context.Context. It is the single seam every layer agrees on:
//
//   - middleware resolves the active org and stores it (WithOrgID),
//   - the postgres OrgScope reads it to set the RLS GUC + restricted role,
//   - tenant repositories read it to filter and stamp rows by org_id.
//
// It is a leaf utility: it imports only the Go standard library, so it can be
// shared by handler/, adapter/ and app/ without violating the dependency rule.
package orgctx

import "context"

// key is an unexported type so the context value cannot collide with keys set by
// other packages.
type key struct{}

// orgIDKey is the single context key under which the active org ID is stored.
var orgIDKey = key{}

// WithOrgID returns a copy of ctx carrying the active organization ID.
func WithOrgID(ctx context.Context, orgID string) context.Context {
	return context.WithValue(ctx, orgIDKey, orgID)
}

// OrgID returns the active organization ID from ctx and whether one was present.
// A missing or empty value reports ok=false.
func OrgID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(orgIDKey).(string)
	if !ok || v == "" {
		return "", false
	}
	return v, true
}
