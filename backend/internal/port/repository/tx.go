package repository

import "context"

// TxManager is the application's transaction abstraction (a Unit-of-Work seam).
//
// It lets a use case run several repository writes atomically without the application
// layer ever importing database/sql or knowing about *sql.Tx. The adapter (postgres)
// begins a real transaction, builds repositories bound to that transaction, and hands
// them to the callback; it commits if the callback returns nil and rolls back fully on
// any error (or panic). Mocks in tests just invoke the callback directly.
//
// Reusability: each multi-write flow gets one focused method here (interface
// segregation — no god "Tx" object). WithTeamTx covers the team-create flow; a new
// flow such as billing would add e.g. WithBillingTx(ctx, fn(subs, invoices) error)
// following the exact same shape. Single-write flows keep using their repository
// directly and need nothing from this interface.
type TxManager interface {
	// WithTeamTx runs fn inside one transaction, passing transaction-scoped team and
	// member repositories. Both repositories' writes commit together or roll back
	// together. Used by the team-create flow so a team and its owner-member are never
	// persisted independently.
	WithTeamTx(ctx context.Context, fn func(teams TeamRepository, members TeamMemberRepository) error) error
}
