package repository

import (
	"context"
	"time"
)

type PasswordReset struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	Used      bool
	CreatedAt time.Time
}

type PasswordResetRepository interface {
	Create(ctx context.Context, pr *PasswordReset) error
	FindByToken(ctx context.Context, token string) (*PasswordReset, error)
	MarkUsed(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) error
}
