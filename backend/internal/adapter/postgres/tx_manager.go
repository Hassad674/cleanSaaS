package postgres

import (
	"context"
	"database/sql"

	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
)

// TxManager implements repository.TxManager on top of *sql.DB. It is the concrete
// Unit-of-Work for the postgres adapter: it begins a transaction, builds tx-scoped
// repositories that write through that transaction, and commits/rolls back atomically.
//
// To add another multi-write flow, add a method here that mirrors WithTeamTx (begin a
// tx via WithTx, build the tx-scoped repos that flow needs, pass them to the callback)
// and a matching method on the repository.TxManager port.
type TxManager struct {
	db *sql.DB
}

// NewTxManager creates a TxManager backed by the connection pool.
func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{db: db}
}

// compile-time check that TxManager satisfies the port.
var _ repository.TxManager = (*TxManager)(nil)

// WithTeamTx runs fn inside a single transaction with transaction-scoped team and
// member repositories. If fn returns an error (or panics), the whole transaction —
// both the team insert and the owner-member insert — is rolled back.
func (m *TxManager) WithTeamTx(ctx context.Context, fn func(teams repository.TeamRepository, members repository.TeamMemberRepository) error) error {
	return WithTx(ctx, m.db, func(tx *sql.Tx) error {
		teams := newTeamRepositoryTx(tx)
		members := newTeamMemberRepositoryTx(tx)
		return fn(teams, members)
	})
}
