package postgres

import (
	"context"
	"database/sql"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO users (email, name, password_hash, avatar_url, role, email_verified, stripe_id, provider, provider_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		 RETURNING id`,
		u.Email, u.Name, u.PasswordHash, u.AvatarURL, u.Role, u.EmailVerified, u.StripeID, u.Provider, u.ProviderID, u.CreatedAt, u.UpdatedAt,
	).Scan(&u.ID)
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	u := &user.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, name, password_hash, avatar_url, role, email_verified, stripe_id, provider, provider_id, created_at, updated_at
		 FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.AvatarURL, &u.Role, &u.EmailVerified, &u.StripeID, &u.Provider, &u.ProviderID, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return u, err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	u := &user.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, name, password_hash, avatar_url, role, email_verified, stripe_id, provider, provider_id, created_at, updated_at
		 FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.AvatarURL, &u.Role, &u.EmailVerified, &u.StripeID, &u.Provider, &u.ProviderID, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return u, err
}

func (r *UserRepository) FindByProvider(ctx context.Context, provider, providerID string) (*user.User, error) {
	u := &user.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, name, password_hash, avatar_url, role, email_verified, stripe_id, provider, provider_id, created_at, updated_at
		 FROM users WHERE provider = $1 AND provider_id = $2`, provider, providerID,
	).Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.AvatarURL, &u.Role, &u.EmailVerified, &u.StripeID, &u.Provider, &u.ProviderID, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return u, err
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET email=$1, name=$2, password_hash=$3, avatar_url=$4, role=$5, email_verified=$6, stripe_id=$7, updated_at=$8
		 WHERE id=$9`,
		u.Email, u.Name, u.PasswordHash, u.AvatarURL, u.Role, u.EmailVerified, u.StripeID, u.UpdatedAt, u.ID,
	)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}

func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]*user.User, int, error) {
	var total int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, email, name, password_hash, avatar_url, role, email_verified, stripe_id, provider, provider_id, created_at, updated_at
		 FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		u := &user.User{}
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.AvatarURL, &u.Role, &u.EmailVerified, &u.StripeID, &u.Provider, &u.ProviderID, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, nil
}
