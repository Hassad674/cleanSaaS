package postgres

import (
	"context"
	"database/sql"

	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
	"github.com/hassad/boilerplateSaaS/backend/pkg/orgctx"
)

// orgFilter returns the active org id from the context as a query argument for the
// repositories' `WHERE org_id = $n` clauses (defense layer 2). When no org is in
// context it returns nil, so the comparison becomes `org_id = NULL` — which is
// never true, yielding zero rows. That makes layer 2 deny-by-default too, matching
// the RLS policy, so a system path that forgot to set an org sees nothing rather
// than everything.
func orgFilter(ctx context.Context) any {
	if id, ok := orgctx.OrgID(ctx); ok {
		return id
	}
	return nil
}

// OrgScope is the concrete org-scoped unit-of-work for the postgres adapter. It
// implements every per-feature scope port (FileScope, ConversationScope, …) by
// opening a transaction that switches to the restricted RLS role and sets the
// active-org GUC (setOrgScope), then handing the callback a transaction-bound
// repository. Row-level security therefore enforces tenant isolation for the whole
// callback, and the short transaction is committed (or rolled back) atomically.
//
// The active org is read from the context. A request with no active org cannot use
// these scopes — setOrgScope rejects a missing/invalid org id — which is the
// intended fail-closed behavior on the tenant request path.
type OrgScope struct {
	db *sql.DB
}

// NewOrgScope creates an org-scoped unit-of-work over the connection pool.
func NewOrgScope(db *sql.DB) *OrgScope {
	return &OrgScope{db: db}
}

// Compile-time checks that OrgScope satisfies every per-feature scope port.
var (
	_ repository.FileScope         = (*OrgScope)(nil)
	_ repository.ConversationScope = (*OrgScope)(nil)
	_ repository.NotificationScope = (*OrgScope)(nil)
	_ repository.SubscriptionScope = (*OrgScope)(nil)
)

func (s *OrgScope) WithOrgFiles(ctx context.Context, fn func(files repository.FileRepository) error) error {
	return s.run(ctx, func(tx *sql.Tx) error { return fn(newFileRepositoryTx(tx)) })
}

func (s *OrgScope) WithOrgConversations(ctx context.Context, fn func(conversations repository.ConversationRepository) error) error {
	return s.run(ctx, func(tx *sql.Tx) error { return fn(newConversationRepositoryTx(tx)) })
}

func (s *OrgScope) WithOrgNotifications(ctx context.Context, fn func(notifications repository.NotificationRepository) error) error {
	return s.run(ctx, func(tx *sql.Tx) error { return fn(newNotificationRepositoryTx(tx)) })
}

func (s *OrgScope) WithOrgSubscriptions(ctx context.Context, fn func(subscriptions repository.SubscriptionRepository) error) error {
	return s.run(ctx, func(tx *sql.Tx) error { return fn(newSubscriptionRepositoryTx(tx)) })
}

// run opens a transaction, applies the RLS role + org GUC for the active org, and
// invokes fn with the transaction. The org id is validated as a UUID and bound as
// a parameter inside setOrgScope — never string-concatenated into SQL.
func (s *OrgScope) run(ctx context.Context, fn func(tx *sql.Tx) error) error {
	orgID, _ := orgctx.OrgID(ctx)
	return WithTx(ctx, s.db, func(tx *sql.Tx) error {
		if err := setOrgScope(ctx, tx, orgID); err != nil {
			return err
		}
		return fn(tx)
	})
}
