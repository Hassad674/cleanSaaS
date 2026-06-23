package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
)

type EmailVerificationRepository struct {
	db *sql.DB
}

func NewEmailVerificationRepository(db *sql.DB) *EmailVerificationRepository {
	return &EmailVerificationRepository{db: db}
}

func (r *EmailVerificationRepository) Create(ctx context.Context, ev *repository.EmailVerification) error {
	query := `INSERT INTO email_verifications (user_id, token, expires_at) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := r.db.QueryRowContext(ctx, query, ev.UserID, ev.Token, ev.ExpiresAt).Scan(&ev.ID, &ev.CreatedAt)
	if err != nil {
		return fmt.Errorf("inserting email verification: %w", err)
	}
	return nil
}

func (r *EmailVerificationRepository) FindByToken(ctx context.Context, token string) (*repository.EmailVerification, error) {
	query := `SELECT id, user_id, token, expires_at, created_at FROM email_verifications WHERE token = $1`
	ev := &repository.EmailVerification{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(&ev.ID, &ev.UserID, &ev.Token, &ev.ExpiresAt, &ev.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("finding email verification by token: %w", err)
	}
	return ev, nil
}

func (r *EmailVerificationRepository) DeleteByUserID(ctx context.Context, userID string) error {
	query := `DELETE FROM email_verifications WHERE user_id = $1`
	if _, err := r.db.ExecContext(ctx, query, userID); err != nil {
		return fmt.Errorf("deleting email verifications by user: %w", err)
	}
	return nil
}

func (r *EmailVerificationRepository) DeleteExpired(ctx context.Context) error {
	ctx, cancel := ctxWithTimeout(ctx, defaultDBTimeout)
	defer cancel()

	query := `DELETE FROM email_verifications WHERE expires_at < NOW()`
	if _, err := r.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("deleting expired email verifications: %w", err)
	}
	return nil
}
