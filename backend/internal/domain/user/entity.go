package user

import (
	"time"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
)

type User struct {
	ID            string
	Email         string
	Name          string
	PasswordHash  string
	AvatarURL     string
	Role          Role
	EmailVerified bool
	StripeID      string
	Provider      string // "email", "google", "github"
	ProviderID    string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func New(email, name, passwordHash string) (*User, error) {
	if email == "" {
		return nil, domain.ErrValidation
	}
	if name == "" {
		return nil, domain.ErrValidation
	}

	return &User{
		Email:        email,
		Name:         name,
		PasswordHash: passwordHash,
		Role:         RoleMember,
		Provider:     "email",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func NewOAuth(email, name, provider, providerID string) *User {
	return &User{
		Email:         email,
		Name:          name,
		Role:          RoleMember,
		EmailVerified: true,
		Provider:      provider,
		ProviderID:    providerID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) VerifyEmail() {
	u.EmailVerified = true
	u.UpdatedAt = time.Now()
}
