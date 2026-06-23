package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
)

type PasswordResetRepository struct {
	db *sql.DB
}

func NewPasswordResetRepository(db *sql.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

func (r *PasswordResetRepository) Create(ctx context.Context, pr *repository.PasswordReset) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO password_resets (user_id, token, expires_at, used, created_at)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		pr.UserID, pr.Token, pr.ExpiresAt, pr.Used, pr.CreatedAt,
	).Scan(&pr.ID)
	if err != nil {
		return fmt.Errorf("creating password reset: %w", err)
	}
	return nil
}

func (r *PasswordResetRepository) FindByToken(ctx context.Context, token string) (*repository.PasswordReset, error) {
	pr := &repository.PasswordReset{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, token, expires_at, used, created_at
		 FROM password_resets WHERE token = $1`, token,
	).Scan(&pr.ID, &pr.UserID, &pr.Token, &pr.ExpiresAt, &pr.Used, &pr.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("finding password reset by token: %w", err)
	}
	return pr, nil
}

func (r *PasswordResetRepository) MarkUsed(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE password_resets SET used = true WHERE id = $1`, id,
	)
	if err != nil {
		return fmt.Errorf("marking password reset as used: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking rows affected: %w", err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *PasswordResetRepository) DeleteExpired(ctx context.Context) error {
	ctx, cancel := ctxWithTimeout(ctx, defaultDBTimeout)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		`DELETE FROM password_resets WHERE expires_at < NOW() OR used = true`,
	)
	if err != nil {
		return fmt.Errorf("deleting expired password resets: %w", err)
	}
	return nil
}
