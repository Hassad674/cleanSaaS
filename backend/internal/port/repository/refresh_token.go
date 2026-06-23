package repository

import (
	"context"
	"time"
)

// RefreshToken is a persisted, hashed refresh token. The raw opaque token is
// never stored — only TokenHash (SHA-256 hex). A token is valid when it exists,
// has not expired, and has not been revoked.
type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

// IsValid reports whether the token can still be used to obtain a new access
// token at time t (not revoked and not expired).
func (rt *RefreshToken) IsValid(t time.Time) bool {
	return rt.RevokedAt == nil && t.Before(rt.ExpiresAt)
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *RefreshToken) error
	FindByHash(ctx context.Context, hash string) (*RefreshToken, error)
	Revoke(ctx context.Context, hash string) error
	RevokeAllForUser(ctx context.Context, userID string) error
	DeleteExpired(ctx context.Context) error
}
