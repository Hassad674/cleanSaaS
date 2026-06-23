package repository

import "context"

// Org-scoped unit-of-work seams. Each tenant feature depends on a small scope
// interface instead of a bare repository. The adapter implements the scope by
// opening a transaction that switches to the restricted RLS role and sets the
// active-org GUC, then handing the callback a transaction-bound repository. Every
// query the callback makes is therefore enforced by PostgreSQL row-level security
// (defense layer 3) in addition to the repository's own WHERE org_id (layer 2).
//
// The app layer never imports database/sql: it only knows "run my work scoped to
// the active org". Tests implement these by invoking the callback directly with a
// plain mock repository, so unit tests stay infrastructure-free.
//
// The active org id is read from the context (see pkg/orgctx) that the request
// path populates; system paths (which carry no org) must not use these scopes.
//
// Interface segregation: one focused scope per feature, mirroring the existing
// per-flow TxManager methods — no god "scope everything" object.

// FileScope runs file-repository work scoped to the active organization.
type FileScope interface {
	WithOrgFiles(ctx context.Context, fn func(files FileRepository) error) error
}

// ConversationScope runs conversation-repository work scoped to the active org.
type ConversationScope interface {
	WithOrgConversations(ctx context.Context, fn func(conversations ConversationRepository) error) error
}

// NotificationScope runs notification-repository work scoped to the active org.
type NotificationScope interface {
	WithOrgNotifications(ctx context.Context, fn func(notifications NotificationRepository) error) error
}

// SubscriptionScope runs subscription-repository work scoped to the active org.
type SubscriptionScope interface {
	WithOrgSubscriptions(ctx context.Context, fn func(subscriptions SubscriptionRepository) error) error
}
