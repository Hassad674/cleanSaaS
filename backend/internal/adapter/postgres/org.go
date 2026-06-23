package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/org"
)

// OrganizationRepository implements repository.OrganizationRepository.
// It holds a DBTX so its writes can run on either the pool or an open transaction
// (used by the signup unit-of-work). Organizations are a system table: they are
// looked up/created on the privileged connection, not under the RLS role.
type OrganizationRepository struct {
	db DBTX
}

// NewOrganizationRepository creates an org repository bound to the pool.
func NewOrganizationRepository(db *sql.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

// newOrganizationRepositoryTx binds the repository to an open transaction.
func newOrganizationRepositoryTx(tx DBTX) *OrganizationRepository {
	return &OrganizationRepository{db: tx}
}

func (r *OrganizationRepository) Create(ctx context.Context, o *org.Organization) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO organizations (name, slug, owner_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		o.Name, o.Slug, o.OwnerID, o.CreatedAt, o.UpdatedAt,
	).Scan(&o.ID)
	if err != nil {
		return fmt.Errorf("creating organization: %w", err)
	}
	return nil
}

func (r *OrganizationRepository) FindByID(ctx context.Context, id string) (*org.Organization, error) {
	o := &org.Organization{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, slug, owner_id, created_at, updated_at FROM organizations WHERE id = $1`, id,
	).Scan(&o.ID, &o.Name, &o.Slug, &o.OwnerID, &o.CreatedAt, &o.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("finding organization by id: %w", err)
	}
	return o, nil
}

func (r *OrganizationRepository) FindBySlug(ctx context.Context, slug string) (*org.Organization, error) {
	o := &org.Organization{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, slug, owner_id, created_at, updated_at FROM organizations WHERE slug = $1`, slug,
	).Scan(&o.ID, &o.Name, &o.Slug, &o.OwnerID, &o.CreatedAt, &o.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("finding organization by slug: %w", err)
	}
	return o, nil
}

func (r *OrganizationRepository) FindDefaultForUser(ctx context.Context, userID string) (*org.Organization, error) {
	o := &org.Organization{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, slug, owner_id, created_at, updated_at
		 FROM organizations WHERE owner_id = $1
		 ORDER BY created_at ASC LIMIT 1`, userID,
	).Scan(&o.ID, &o.Name, &o.Slug, &o.OwnerID, &o.CreatedAt, &o.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("finding default organization for user: %w", err)
	}
	return o, nil
}

// OrganizationMemberRepository implements repository.OrganizationMemberRepository.
type OrganizationMemberRepository struct {
	db DBTX
}

// NewOrganizationMemberRepository creates a member repository bound to the pool.
func NewOrganizationMemberRepository(db *sql.DB) *OrganizationMemberRepository {
	return &OrganizationMemberRepository{db: db}
}

// newOrganizationMemberRepositoryTx binds the repository to an open transaction.
func newOrganizationMemberRepositoryTx(tx DBTX) *OrganizationMemberRepository {
	return &OrganizationMemberRepository{db: tx}
}

func (r *OrganizationMemberRepository) Add(ctx context.Context, m *org.Member) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO organization_members (org_id, user_id, role, created_at)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (org_id, user_id) DO NOTHING`,
		m.OrgID, m.UserID, string(m.Role), m.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("adding organization member: %w", err)
	}
	return nil
}

func (r *OrganizationMemberRepository) FindByOrgAndUser(ctx context.Context, orgID, userID string) (*org.Member, error) {
	m := &org.Member{}
	var role string
	err := r.db.QueryRowContext(ctx,
		`SELECT org_id, user_id, role, created_at FROM organization_members WHERE org_id = $1 AND user_id = $2`,
		orgID, userID,
	).Scan(&m.OrgID, &m.UserID, &role, &m.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("finding organization member: %w", err)
	}
	m.Role = org.Role(role)
	return m, nil
}

func (r *OrganizationMemberRepository) IsMember(ctx context.Context, orgID, userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM organization_members WHERE org_id = $1 AND user_id = $2)`,
		orgID, userID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking organization membership: %w", err)
	}
	return exists, nil
}
