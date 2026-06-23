package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

// DBTX is the subset of methods shared by *sql.DB and *sql.Tx. Repositories hold a
// DBTX instead of a concrete *sql.DB, so the exact same repository code runs whether
// it was given the pooled connection (*sql.DB) or a transaction handle (*sql.Tx).
// This is what lets a tx-scoped repository join a caller's transaction transparently.
//
// Both *sql.DB and *sql.Tx already implement these three methods with these exact
// signatures, so passing either one to a NewXRepository(db *sql.DB) constructor (or
// to newXRepositoryTx(tx)) just works — no wrapper needed.
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// WithTx runs fn inside a single database transaction: it begins the tx, commits on
// success, and rolls back on any error or panic. On panic it rolls back and re-panics
// so the failure is never silently swallowed.
//
// This is the low-level primitive used by TxManager. Most callers should depend on the
// repository.TxManager port instead of calling this directly, so the application layer
// stays free of database/sql.
func WithTx(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			// A panic must not commit a half-applied transaction.
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("rollback failed: %w (original: %v)", rbErr, err)
			}
			return
		}
		if cmErr := tx.Commit(); cmErr != nil {
			err = fmt.Errorf("commit tx: %w", cmErr)
		}
	}()

	return fn(tx)
}
