//go:build integration

package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/user"
)

// TestUserRepository_RoundTrip exercises the core user CRUD path against the real
// schema: create (RETURNING id), find by id, find by email, update, and the
// not-found mapping. Cleanup is handled by newUser's registered delete.
func TestUserRepository_RoundTrip(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	repo := NewUserRepository(db)

	u := newUser(ctx, t, db)

	// FindByID returns the persisted row with the generated UUID.
	got, err := repo.FindByID(ctx, u.ID)
	require.NoError(t, err)
	assert.Equal(t, u.ID, got.ID)
	assert.Equal(t, u.Email, got.Email)
	assert.Equal(t, "hash", got.PasswordHash)
	assert.Equal(t, user.RoleMember, got.Role)

	// FindByEmail resolves the same row.
	byEmail, err := repo.FindByEmail(ctx, u.Email)
	require.NoError(t, err)
	assert.Equal(t, u.ID, byEmail.ID)

	// Update mutates persisted columns.
	u.Name = "Renamed"
	u.EmailVerified = true
	u.StripeID = "cus_" + uniqueTag()
	u.UpdatedAt = time.Now()
	require.NoError(t, repo.Update(ctx, u))

	updated, err := repo.FindByID(ctx, u.ID)
	require.NoError(t, err)
	assert.Equal(t, "Renamed", updated.Name)
	assert.True(t, updated.EmailVerified)
	assert.Equal(t, u.StripeID, updated.StripeID)

	// FindByStripeID resolves via the just-set stripe id.
	byStripe, err := repo.FindByStripeID(ctx, u.StripeID)
	require.NoError(t, err)
	assert.Equal(t, u.ID, byStripe.ID)
}

// TestUserRepository_NotFound proves the sql.ErrNoRows -> domain.ErrNotFound
// mapping for a missing user, using a random UUID that cannot exist.
func TestUserRepository_NotFound(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	repo := NewUserRepository(db)

	_, err := repo.FindByID(ctx, uniqueTag())
	assert.ErrorIs(t, err, domain.ErrNotFound)

	_, err = repo.FindByEmail(ctx, "nobody-"+uniqueTag()+"@itest.test")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

// TestUserRepository_Delete proves a created user is gone after Delete and that
// a subsequent lookup reports not-found.
func TestUserRepository_Delete(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	repo := NewUserRepository(db)

	tag := uniqueTag()
	u, err := user.New("itest-del-"+tag+"@itest.test", "ToDelete", "hash")
	require.NoError(t, err)
	require.NoError(t, repo.Create(ctx, u))

	require.NoError(t, repo.Delete(ctx, u.ID))

	_, err = repo.FindByID(ctx, u.ID)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

// TestUserRepository_ListAndSearch_Pagination verifies the COUNT(*) + ORDER BY +
// LIMIT/OFFSET path returns the inserted users and a total >= the number created,
// and that Search narrows by an ILIKE pattern unique to this run.
func TestUserRepository_ListAndSearch_Pagination(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	repo := NewUserRepository(db)

	// A run-unique marker placed in the name so Search matches only our rows.
	marker := "ZZ" + uniqueTag()
	for i := 0; i < 3; i++ {
		tag := uniqueTag()
		u, err := user.New("itest-list-"+tag+"@itest.test", marker+" user", "hash")
		require.NoError(t, err)
		require.NoError(t, repo.Create(ctx, u))
		id := u.ID
		t.Cleanup(func() { _ = repo.Delete(context.Background(), id) })
	}

	// List returns a page and a total that counts at least our 3 new rows.
	page, total, err := repo.List(ctx, 0, 2)
	require.NoError(t, err)
	assert.Len(t, page, 2, "LIMIT 2 yields a 2-row page")
	assert.GreaterOrEqual(t, total, 3, "total counts all users including the 3 just created")

	// Search by the run-unique marker finds exactly the 3 created rows.
	found, searchTotal, err := repo.Search(ctx, marker, 0, 50)
	require.NoError(t, err)
	assert.Equal(t, 3, searchTotal, "exactly the 3 marker users match")
	assert.Len(t, found, 3)

	// Count returns a positive total (>= our rows + seeded users).
	count, err := repo.Count(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, 3)
}
