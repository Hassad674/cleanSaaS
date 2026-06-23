package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// roleAppUser is the NOBYPASSRLS database role the request path switches to so
// PostgreSQL row-level security is actually enforced. Migrations and seed keep
// using the privileged connection role (which owns the tables and bypasses RLS) —
// only the per-request org scope switches to this restricted role. It is a fixed
// constant (never user input), so using it in SET LOCAL ROLE is safe.
const roleAppUser = "app_user"

// setOrgScope applies the row-level-security guardrail to an open transaction:
// switch to the restricted role and set the active-org GUC. SET LOCAL confines
// both to this transaction, so a pooled connection can never carry a leftover role
// or org into the next request. Under app_user (NOBYPASSRLS) the org_isolation
// policies apply — this is defense layer 3.
//
// The org id is validated as a UUID and passed to set_config as a BOUND PARAMETER,
// so it is never string-concatenated into SQL. An empty/invalid org id fails
// closed: the request path cannot proceed without a valid active organization.
func setOrgScope(ctx context.Context, tx *sql.Tx, orgID string) error {
	if _, err := uuid.Parse(orgID); err != nil {
		return fmt.Errorf("invalid org id for RLS scope: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "SET LOCAL ROLE "+roleAppUser); err != nil {
		return fmt.Errorf("setting RLS role: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_org_id', $1, true)", orgID); err != nil {
		return fmt.Errorf("setting org GUC: %w", err)
	}
	return nil
}
