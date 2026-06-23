package user

import (
	"time"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
)

type User struct {
	ID string
	// Email is stored as a normalized string at the repository boundary. It is
	// validated and normalized through the Email value object in New/NewOAuth;
	// use EmailVO() to obtain the typed value object when an invariant needs it.
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
	addr, err := NewEmail(email)
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, domain.ErrValidation
	}

	return &User{
		Email:        addr.String(),
		Name:         name,
		PasswordHash: passwordHash,
		Role:         RoleMember,
		Provider:     "email",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func NewOAuth(email, name, provider, providerID string) *User {
	// OAuth providers supply a vetted email; we still normalize it through the
	// value object so casing/whitespace match the email-signup path. A provider
	// that somehow returns a malformed address falls back to the raw value
	// rather than failing the login (the address is trusted upstream).
	normalized := email
	if addr, err := NewEmail(email); err == nil {
		normalized = addr.String()
	}
	return &User{
		Email:         normalized,
		Name:          name,
		Role:          RoleMember,
		EmailVerified: true,
		Provider:      provider,
		ProviderID:    providerID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// EmailVO returns the user's email as the Email value object. It re-parses the
// stored string (which originated from NewEmail) so callers that need the typed
// invariant — rather than the raw repository string — can obtain it.
func (u *User) EmailVO() (Email, error) {
	return NewEmail(u.Email)
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) VerifyEmail() {
	u.EmailVerified = true
	u.UpdatedAt = time.Now()
}
