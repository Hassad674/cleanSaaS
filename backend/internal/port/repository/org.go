package repository

import (
	"context"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain/org"
)

// OrganizationRepository persists organizations (the tenant aggregate).
//
// Organizations are a SYSTEM-managed table: they are created/looked up on the
// signup and login paths, which run as the privileged (RLS-bypassing) role, so
// these methods are NOT org-scoped by RLS themselves. They are the source of the
// org_id that every tenant-scoped query is then keyed on.
type OrganizationRepository interface {
	Create(ctx context.Context, o *org.Organization) error
	FindByID(ctx context.Context, id string) (*org.Organization, error)
	FindBySlug(ctx context.Context, slug string) (*org.Organization, error)
	// FindDefaultForUser returns the user's primary organization — the one they
	// own (their personal org), used as the active org when a request carries no
	// explicit org selection.
	FindDefaultForUser(ctx context.Context, userID string) (*org.Organization, error)
}

// OrganizationMemberRepository persists organization memberships.
type OrganizationMemberRepository interface {
	Add(ctx context.Context, m *org.Member) error
	FindByOrgAndUser(ctx context.Context, orgID, userID string) (*org.Member, error)
	// IsMember reports whether userID belongs to orgID. Used by middleware to
	// authorize an explicitly-requested active org before trusting it.
	IsMember(ctx context.Context, orgID, userID string) (bool, error)
}
