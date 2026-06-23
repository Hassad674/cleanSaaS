//go:build integration

// Cross-tenant isolation test for PostgreSQL row-level security.
//
// Run against the live local database (port 5433) AFTER applying migrations:
//
//	go run cmd/migrate/main.go up
//	go test -tags=integration ./internal/adapter/postgres/ -run RLS -v
//
// It proves the LAST line of defense (layer 3): under the restricted app_user role
// with app.current_org_id set to org A, a session sees ONLY org A's rows, cannot
// read org B's rows, and is rejected by the WITH CHECK policy when it tries to
// write a row belonging to org B. This holds regardless of any WHERE clause — it is
// the database itself enforcing the tenant boundary.
package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDSN = "postgres://postgres:postgres@localhost:5433/cleansaas?sslmode=disable"

// asAppUser runs fn inside a transaction scoped to orgID under the restricted
// app_user role — exactly what the request path does in production (setOrgScope).
func asAppUser(ctx context.Context, t *testing.T, db *sql.DB, orgID string, fn func(tx *sql.Tx)) {
	t.Helper()
	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer func() { _ = tx.Rollback() }()
	require.NoError(t, setOrgScope(ctx, tx, orgID))
	fn(tx)
}

func TestRLS_CrossTenantIsolation(t *testing.T) {
	ctx := context.Background()
	db, err := sql.Open("postgres", testDSN)
	require.NoError(t, err, "open db")
	// Close the pool LAST: t.Cleanup is LIFO, so registering Close first means it
	// runs after the data-deletion cleanup (registered later) — otherwise the pool
	// would be closed before cleanup and the deletes would silently no-op.
	t.Cleanup(func() { _ = db.Close() })
	require.NoError(t, db.Ping(), "the integration test needs the local DB on :5433 with migrations applied")

	// --- Arrange (as the privileged/superuser role, which bypasses RLS) ---
	// Two users (org owners), two orgs, one file in each org.
	userA := mustCreateUser(ctx, t, db, "rls-a")
	userB := mustCreateUser(ctx, t, db, "rls-b")
	orgA := mustCreateOrg(ctx, t, db, "RLS Org A", userA)
	orgB := mustCreateOrg(ctx, t, db, "RLS Org B", userB)
	t.Cleanup(func() {
		// Deleting the orgs cascades to their files; deleting users cascades to orgs.
		_, _ = db.ExecContext(ctx, `DELETE FROM users WHERE id = ANY($1)`, sqlArray(userA, userB))
	})

	fileA := mustCreateFile(ctx, t, db, orgA, userA, "a.txt", "orgA/a.txt")
	fileB := mustCreateFile(ctx, t, db, orgB, userB, "b.txt", "orgB/b.txt")

	// Sanity: as superuser (RLS bypassed) BOTH files exist.
	assert.Equal(t, 2, countFilesSuper(ctx, t, db, fileA, fileB), "both files exist for the privileged role")

	// --- Act + Assert: under app_user scoped to org A ---
	asAppUser(ctx, t, db, orgA, func(tx *sql.Tx) {
		// 1. SELECT returns ONLY org A's file.
		rows, err := tx.QueryContext(ctx, `SELECT id, org_id FROM files WHERE id = ANY($1)`, sqlArray(fileA, fileB))
		require.NoError(t, err)
		defer rows.Close()
		var seen []string
		for rows.Next() {
			var id, org string
			require.NoError(t, rows.Scan(&id, &org))
			seen = append(seen, id)
			assert.Equal(t, orgA, org, "only org A rows must be visible")
		}
		require.NoError(t, rows.Err())
		assert.Equal(t, []string{fileA}, seen, "exactly org A's file is visible, org B's is hidden")

		// 2. Reading org B's file by id directly returns NO rows (not an error).
		var got string
		err = tx.QueryRowContext(ctx, `SELECT id FROM files WHERE id = $1`, fileB).Scan(&got)
		assert.ErrorIs(t, err, sql.ErrNoRows, "org B's row is invisible under org A's scope")

		// 3. INSERT targeting org B is REJECTED by the WITH CHECK policy.
		_, err = tx.ExecContext(ctx,
			`INSERT INTO files (org_id, user_id, name, key, size_bytes, content_type, url)
			 VALUES ($1, $2, 'evil.txt', $3, 1, 'text/plain', '')`,
			orgB, userA, "orgB/evil-"+uuid.NewString(),
		)
		require.Error(t, err, "inserting a row for another tenant must be rejected")
		assert.True(t, isRLSViolation(err), "expected an RLS WITH CHECK violation, got: %v", err)
	})

	// A fresh transaction (the previous one is poisoned after the failed insert).
	asAppUser(ctx, t, db, orgA, func(tx *sql.Tx) {
		// 4. UPDATE targeting org B's row affects ZERO rows (invisible under USING).
		res, err := tx.ExecContext(ctx, `UPDATE files SET name = 'hacked' WHERE id = $1`, fileB)
		require.NoError(t, err)
		affected, _ := res.RowsAffected()
		assert.Equal(t, int64(0), affected, "cannot update another tenant's row")

		// 5. INSERT for the CORRECT org (org A) succeeds.
		_, err = tx.ExecContext(ctx,
			`INSERT INTO files (org_id, user_id, name, key, size_bytes, content_type, url)
			 VALUES ($1, $2, 'ok.txt', $3, 1, 'text/plain', '')`,
			orgA, userA, "orgA/ok-"+uuid.NewString(),
		)
		assert.NoError(t, err, "inserting a row for the active tenant must succeed")
	})

	// 6. With NO org GUC set, app_user sees NOTHING (deny-by-default).
	func() {
		tx, err := db.BeginTx(ctx, nil)
		require.NoError(t, err)
		defer func() { _ = tx.Rollback() }()
		_, err = tx.ExecContext(ctx, "SET LOCAL ROLE "+roleAppUser)
		require.NoError(t, err)
		var n int
		require.NoError(t, tx.QueryRowContext(ctx, `SELECT count(*) FROM files`).Scan(&n))
		assert.Equal(t, 0, n, "with no app.current_org_id, RLS denies all rows by default")
	}()

	// 7. Org B's file is STILL intact (never updated/hacked from org A's session).
	var stillB string
	require.NoError(t, db.QueryRowContext(ctx, `SELECT name FROM files WHERE id = $1`, fileB).Scan(&stillB))
	assert.Equal(t, "b.txt", stillB, "org B's data was never mutated by org A")
}

// TestRLS_MessagesTransitiveIsolation proves the messages table — which has no
// org_id of its own — is isolated transitively through its parent conversation.
// Under org A's scope a message can be inserted into and read from org A's
// conversation, but org B's conversation and its messages are invisible.
func TestRLS_MessagesTransitiveIsolation(t *testing.T) {
	ctx := context.Background()
	db, err := sql.Open("postgres", testDSN)
	require.NoError(t, err)
	// Close LAST (LIFO) so the data-deletion cleanup runs while the pool is open.
	t.Cleanup(func() { _ = db.Close() })
	require.NoError(t, db.Ping(), "the integration test needs the local DB on :5433 with migrations applied")

	userA := mustCreateUser(ctx, t, db, "rlsm-a")
	userB := mustCreateUser(ctx, t, db, "rlsm-b")
	orgA := mustCreateOrg(ctx, t, db, "RLS-M Org A", userA)
	orgB := mustCreateOrg(ctx, t, db, "RLS-M Org B", userB)
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, `DELETE FROM users WHERE id = ANY($1)`, sqlArray(userA, userB))
	})

	convA := mustCreateConversation(ctx, t, db, orgA, userA)
	convB := mustCreateConversation(ctx, t, db, orgB, userB)
	// Seed one message in org B's conversation (as the privileged role).
	_, err = db.ExecContext(ctx, `INSERT INTO messages (conversation_id, role, content) VALUES ($1, 'user', 'secret B')`, convB)
	require.NoError(t, err)

	// Insert + read within ONE org-A-scoped transaction (the test helper rolls the
	// transaction back at the end, so the insert is observed only within it — that
	// is exactly the lifetime of a real request's org scope).
	asAppUser(ctx, t, db, orgA, func(tx *sql.Tx) {
		// Insert a message into org A's own conversation — allowed, and visible.
		_, err := tx.ExecContext(ctx, `INSERT INTO messages (conversation_id, role, content) VALUES ($1, 'user', 'hello A')`, convA)
		require.NoError(t, err, "message into own-org conversation must be allowed")

		var n int
		require.NoError(t, tx.QueryRowContext(ctx, `SELECT count(*) FROM messages WHERE conversation_id = $1`, convA).Scan(&n))
		assert.Equal(t, 1, n, "org A's own message is visible within its scope")

		// Org B's messages are invisible under org A's scope (transitive isolation).
		require.NoError(t, tx.QueryRowContext(ctx, `SELECT count(*) FROM messages WHERE conversation_id = $1`, convB).Scan(&n))
		assert.Equal(t, 0, n, "org B's messages are invisible under org A's scope")
	})

	// A fresh transaction: inserting into org B's conversation is rejected by the
	// WITH CHECK policy (org B's conversation is not visible under org A's scope).
	asAppUser(ctx, t, db, orgA, func(tx *sql.Tx) {
		_, err := tx.ExecContext(ctx, `INSERT INTO messages (conversation_id, role, content) VALUES ($1, 'user', 'evil')`, convB)
		assert.True(t, isRLSViolation(err), "message into another tenant's conversation must be rejected, got: %v", err)
	})
}

// --- helpers (privileged role) ---

func mustCreateConversation(ctx context.Context, t *testing.T, db *sql.DB, orgID, userID string) string {
	t.Helper()
	var id string
	err := db.QueryRowContext(ctx,
		`INSERT INTO conversations (org_id, user_id, title) VALUES ($1, $2, 'c') RETURNING id`,
		orgID, userID,
	).Scan(&id)
	require.NoError(t, err)
	return id
}

func mustCreateUser(ctx context.Context, t *testing.T, db *sql.DB, tag string) string {
	t.Helper()
	var id string
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, name, password_hash, role, email_verified)
		 VALUES ($1, $2, '', 'user', true) RETURNING id`,
		tag+"-"+uuid.NewString()+"@rls.test", "RLS "+tag,
	).Scan(&id)
	require.NoError(t, err)
	return id
}

func mustCreateOrg(ctx context.Context, t *testing.T, db *sql.DB, name, ownerID string) string {
	t.Helper()
	var id string
	err := db.QueryRowContext(ctx,
		`INSERT INTO organizations (name, slug, owner_id) VALUES ($1, $2, $3) RETURNING id`,
		name, "rls-"+uuid.NewString(), ownerID,
	).Scan(&id)
	require.NoError(t, err)
	return id
}

func mustCreateFile(ctx context.Context, t *testing.T, db *sql.DB, orgID, userID, name, key string) string {
	t.Helper()
	var id string
	err := db.QueryRowContext(ctx,
		`INSERT INTO files (org_id, user_id, name, key, size_bytes, content_type, url)
		 VALUES ($1, $2, $3, $4, 1, 'text/plain', '') RETURNING id`,
		orgID, userID, name, key+"-"+uuid.NewString(),
	).Scan(&id)
	require.NoError(t, err)
	return id
}

func countFilesSuper(ctx context.Context, t *testing.T, db *sql.DB, ids ...string) int {
	t.Helper()
	var n int
	require.NoError(t, db.QueryRowContext(ctx, `SELECT count(*) FROM files WHERE id = ANY($1)`, sqlArray(ids...)).Scan(&n))
	return n
}

// isRLSViolation reports whether err is a PostgreSQL row-level-security policy
// violation (SQLSTATE 42501 / message mentioning row-level security).
func isRLSViolation(err error) bool {
	if err == nil {
		return false
	}
	return containsAny(err.Error(), "row-level security", "violates row-level security", "42501")
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if len(sub) > 0 && indexOf(s, sub) >= 0 {
			return true
		}
	}
	return false
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

// sqlArray builds a Postgres array parameter for `= ANY($n)` queries from string
// ids, using lib/pq's array encoder so the driver binds a real text[] (a plain Go
// string would not coerce and the query would silently match nothing).
func sqlArray(ids ...string) any {
	return pq.Array(ids)
}
