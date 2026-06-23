//go:build integration

// Shared scaffolding for the postgres adapter integration suite.
//
// All tests in this package are gated behind the `integration` build tag and run
// against the LIVE local PostgreSQL (port 5433) with migrations applied:
//
//	go run cmd/migrate/main.go up
//	go test -tags=integration ./internal/adapter/postgres/... -count=1
//
// The suite degrades gracefully: if DATABASE_URL / the DB is unreachable, each
// test t.Skip()s rather than failing, so a tag-gated CI lane without a database
// still passes. Every test creates its OWN data with unique identifiers and
// cleans up after itself (defer/ t.Cleanup deletes, or an org-scoped transaction
// that rolls back), so the suite is fully re-runnable and never touches the
// seeded admin/demo/plan/blog rows.
package postgres

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain/org"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/user"
	"github.com/hassad/boilerplateSaaS/backend/pkg/orgctx"
)

// dsn resolves the integration database connection string from DATABASE_URL,
// falling back to the same local default the RLS test uses (testDSN). Keeping the
// env override means the suite honors whatever DB the developer/CI points at.
func dsn() string {
	if v := os.Getenv("DATABASE_URL"); v != "" {
		return v
	}
	return testDSN
}

// openTestDB opens the connection pool and pings it. If the DB is unreachable it
// SKIPS the test (graceful degradation for a DB-less CI lane) rather than failing.
// The pool is registered for Close LAST (t.Cleanup is LIFO) so any data-deletion
// cleanups registered later still run while the pool is open.
func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("postgres", dsn())
	if err != nil {
		t.Skipf("skipping integration test: cannot open db: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		t.Skipf("skipping integration test: local DB on :5433 with migrations applied is required: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

// uniqueTag returns a collision-proof token for building unique emails, slugs and
// keys, so re-running the suite never trips a UNIQUE constraint and rows from one
// run never alias another.
func uniqueTag() string {
	return uuid.NewString()
}

// newUser creates a user via the real UserRepository (privileged pool) and
// registers its deletion. Deleting the user cascades to every dependent row
// (orgs, files, subscriptions, …), which keeps cleanup simple and total.
func newUser(ctx context.Context, t *testing.T, db *sql.DB) *user.User {
	t.Helper()
	repo := NewUserRepository(db)
	tag := uniqueTag()
	u, err := user.New("itest-"+tag+"@itest.test", "Integration "+tag, "hash")
	require.NoError(t, err)
	require.NoError(t, repo.Create(ctx, u))
	require.NotEmpty(t, u.ID)
	t.Cleanup(func() { _ = repo.Delete(context.Background(), u.ID) })
	return u
}

// newOrg creates an organization owned by ownerID via the real
// OrganizationRepository (privileged pool, the signup/system path). Cleanup is
// covered transitively by the owning user's cascade, but we also register an
// explicit delete so an org created for a seeded/other owner never leaks.
func newOrg(ctx context.Context, t *testing.T, db *sql.DB, ownerID string) *org.Organization {
	t.Helper()
	repo := NewOrganizationRepository(db)
	tag := uniqueTag()
	o := &org.Organization{
		Name:      "Integration Org " + tag,
		Slug:      "itest-" + tag,
		OwnerID:   ownerID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, repo.Create(ctx, o))
	require.NotEmpty(t, o.ID)
	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(), `DELETE FROM organizations WHERE id = $1`, o.ID)
	})
	return o
}

// orgScopedCtx returns a context carrying orgID, exactly as the request path's
// org middleware would populate it (pkg/orgctx). The org-scoped repositories read
// the active org from here to filter/stamp rows, and OrgScope reads it to set the
// RLS GUC + restricted role.
func orgScopedCtx(orgID string) context.Context {
	return orgctx.WithOrgID(context.Background(), orgID)
}
