//go:build integration

package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/org"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/user"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
)

// TestOrganizationRepository_RoundTrip exercises the org aggregate on the
// privileged (system) path: create, find by id, find by slug, and
// FindDefaultForUser (the owner's earliest org).
func TestOrganizationRepository_RoundTrip(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	repo := NewOrganizationRepository(db)

	owner := newUser(ctx, t, db)
	o := newOrg(ctx, t, db, owner.ID)

	byID, err := repo.FindByID(ctx, o.ID)
	require.NoError(t, err)
	assert.Equal(t, o.Slug, byID.Slug)
	assert.Equal(t, owner.ID, byID.OwnerID)

	bySlug, err := repo.FindBySlug(ctx, o.Slug)
	require.NoError(t, err)
	assert.Equal(t, o.ID, bySlug.ID)

	// FindDefaultForUser returns the earliest-created org owned by the user.
	def, err := repo.FindDefaultForUser(ctx, owner.ID)
	require.NoError(t, err)
	assert.Equal(t, o.ID, def.ID)
}

// TestOrganizationRepository_NotFound proves the not-found mapping for missing
// org lookups by id and slug.
func TestOrganizationRepository_NotFound(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	repo := NewOrganizationRepository(db)

	_, err := repo.FindByID(ctx, uniqueTag())
	assert.ErrorIs(t, err, domain.ErrNotFound)

	_, err = repo.FindBySlug(ctx, "missing-"+uniqueTag())
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

// TestOrganizationMemberRepository_AddAndQuery exercises membership: Add (with
// ON CONFLICT DO NOTHING idempotency), FindByOrgAndUser, and IsMember for both a
// member and a non-member.
func TestOrganizationMemberRepository_AddAndQuery(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	members := NewOrganizationMemberRepository(db)

	owner := newUser(ctx, t, db)
	o := newOrg(ctx, t, db, owner.ID)
	member := newUser(ctx, t, db)

	m := &org.Member{
		OrgID:     o.ID,
		UserID:    member.ID,
		Role:      org.RoleMember,
		CreatedAt: time.Now(),
	}
	require.NoError(t, members.Add(ctx, m))

	// ON CONFLICT (org_id, user_id) DO NOTHING: re-adding the same membership is
	// a harmless no-op, not an error — proving the idempotent insert.
	require.NoError(t, members.Add(ctx, m), "duplicate Add must be a silent no-op")

	got, err := members.FindByOrgAndUser(ctx, o.ID, member.ID)
	require.NoError(t, err)
	assert.Equal(t, org.RoleMember, got.Role)

	isMember, err := members.IsMember(ctx, o.ID, member.ID)
	require.NoError(t, err)
	assert.True(t, isMember)

	// A user who was never added is not a member.
	stranger := newUser(ctx, t, db)
	isMember, err = members.IsMember(ctx, o.ID, stranger.ID)
	require.NoError(t, err)
	assert.False(t, isMember)

	// Missing membership maps to not-found.
	_, err = members.FindByOrgAndUser(ctx, o.ID, stranger.ID)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

// TestTxManager_WithSignupTx_AtomicCommit proves the signup unit-of-work commits
// the user, their org, and the owner membership together, and that they are all
// visible afterwards. Cleanup deletes the user (cascades to org + membership).
func TestTxManager_WithSignupTx_AtomicCommit(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	txm := NewTxManager(db)

	tag := uniqueTag()
	var createdUserID, createdOrgID string

	err := txm.WithSignupTx(ctx, func(users repository.UserRepository, orgs repository.OrganizationRepository, mems repository.OrganizationMemberRepository) error {
		u, e := user.New("itest-signup-"+tag+"@itest.test", "Signup "+tag, "hash")
		if e != nil {
			return e
		}
		if e := users.Create(ctx, u); e != nil {
			return e
		}
		createdUserID = u.ID

		o := &org.Organization{Name: "Signup Org " + tag, Slug: "itest-signup-" + tag, OwnerID: u.ID, CreatedAt: time.Now(), UpdatedAt: time.Now()}
		if e := orgs.Create(ctx, o); e != nil {
			return e
		}
		createdOrgID = o.ID

		return mems.Add(ctx, &org.Member{OrgID: o.ID, UserID: u.ID, Role: org.RoleOwner, CreatedAt: time.Now()})
	})
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = db.ExecContext(context.Background(), `DELETE FROM users WHERE id = $1`, createdUserID) })

	// All three rows committed together and are now readable on the pool.
	u, err := NewUserRepository(db).FindByID(ctx, createdUserID)
	require.NoError(t, err)
	assert.Equal(t, createdUserID, u.ID)

	isMember, err := NewOrganizationMemberRepository(db).IsMember(ctx, createdOrgID, createdUserID)
	require.NoError(t, err)
	assert.True(t, isMember, "owner membership committed atomically with the user and org")
}
