package repository

import (
	"context"
	"time"
)

type EmailVerification struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type EmailVerificationRepository interface {
	Create(ctx context.Context, ev *EmailVerification) error
	FindByToken(ctx context.Context, token string) (*EmailVerification, error)
	DeleteByUserID(ctx context.Context, userID string) error
	DeleteExpired(ctx context.Context) error
}
