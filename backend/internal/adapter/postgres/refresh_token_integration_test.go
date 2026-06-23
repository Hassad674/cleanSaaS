//go:build integration

package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
)

// TestRefreshTokenRepository_RoundTrip exercises the full token lifecycle:
// create (RETURNING id, created_at), find-by-hash, revoke (with RowsAffected
// guard), and re-revoke reporting not-found.
func TestRefreshTokenRepository_RoundTrip(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	repo := NewRefreshTokenRepository(db)

	u := newUser(ctx, t, db)
	hash := "rt-" + uniqueTag()
	rt := &repository.RefreshToken{
		UserID:    u.ID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	require.NoError(t, repo.Create(ctx, rt))
	require.NotEmpty(t, rt.ID, "Create populates id via RETURNING")
	require.False(t, rt.CreatedAt.IsZero(), "Create populates created_at via RETURNING")

	// FindByHash returns a valid (not revoked, not expired) token.
	got, err := repo.FindByHash(ctx, hash)
	require.NoError(t, err)
	assert.Equal(t, rt.ID, got.ID)
	assert.True(t, got.IsValid(time.Now()), "freshly created token is valid")
	assert.Nil(t, got.RevokedAt)

	// Revoke sets revoked_at; the token is then invalid.
	require.NoError(t, repo.Revoke(ctx, hash))
	revoked, err := repo.FindByHash(ctx, hash)
	require.NoError(t, err)
	require.NotNil(t, revoked.RevokedAt, "revoked_at is set after Revoke")
	assert.False(t, revoked.IsValid(time.Now()), "a revoked token is no longer valid")

	// Re-revoking an already-revoked token affects zero rows -> not-found.
	err = repo.Revoke(ctx, hash)
	assert.ErrorIs(t, err, domain.ErrNotFound, "revoking an already-revoked token reports not-found via RowsAffected==0")

	// FindByHash on a non-existent hash maps to not-found.
	_, err = repo.FindByHash(ctx, "nope-"+uniqueTag())
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

// TestRefreshTokenRepository_RevokeAllForUser proves that revoking all tokens for
// a user invalidates every still-active token and is a safe no-op when re-run.
func TestRefreshTokenRepository_RevokeAllForUser(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	repo := NewRefreshTokenRepository(db)

	u := newUser(ctx, t, db)
	var hashes []string
	for i := 0; i < 3; i++ {
		h := "rt-all-" + uniqueTag()
		hashes = append(hashes, h)
		require.NoError(t, repo.Create(ctx, &repository.RefreshToken{
			UserID:    u.ID,
			TokenHash: h,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}))
	}

	require.NoError(t, repo.RevokeAllForUser(ctx, u.ID))
	for _, h := range hashes {
		got, err := repo.FindByHash(ctx, h)
		require.NoError(t, err)
		assert.NotNil(t, got.RevokedAt, "every user token is revoked")
	}

	// RevokeAll with nothing left to revoke is a no-op, not an error.
	require.NoError(t, repo.RevokeAllForUser(ctx, u.ID))
}

// TestRefreshTokenRepository_DeleteExpired proves the cleanup query removes
// expired and revoked tokens while leaving a still-valid token untouched. Only
// rows owned by this test's user are inspected, so seeded data is irrelevant.
func TestRefreshTokenRepository_DeleteExpired(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	repo := NewRefreshTokenRepository(db)

	u := newUser(ctx, t, db)

	// A valid token (future expiry, not revoked) — must survive.
	validHash := "rt-valid-" + uniqueTag()
	require.NoError(t, repo.Create(ctx, &repository.RefreshToken{
		UserID:    u.ID,
		TokenHash: validHash,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}))

	// An expired token — must be deleted.
	expiredHash := "rt-expired-" + uniqueTag()
	require.NoError(t, repo.Create(ctx, &repository.RefreshToken{
		UserID:    u.ID,
		TokenHash: expiredHash,
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}))

	// A revoked-but-not-expired token — also deleted (the query removes revoked).
	revokedHash := "rt-revoked-" + uniqueTag()
	require.NoError(t, repo.Create(ctx, &repository.RefreshToken{
		UserID:    u.ID,
		TokenHash: revokedHash,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}))
	require.NoError(t, repo.Revoke(ctx, revokedHash))

	require.NoError(t, repo.DeleteExpired(ctx))

	_, err := repo.FindByHash(ctx, expiredHash)
	assert.ErrorIs(t, err, domain.ErrNotFound, "expired token was deleted")
	_, err = repo.FindByHash(ctx, revokedHash)
	assert.ErrorIs(t, err, domain.ErrNotFound, "revoked token was deleted")

	stillValid, err := repo.FindByHash(ctx, validHash)
	require.NoError(t, err, "valid token survived cleanup")
	assert.True(t, stillValid.IsValid(time.Now()))
}
