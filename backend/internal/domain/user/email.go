package user

import (
	"net/mail"
	"strings"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
)

// Email is a value object representing a validated, normalized email address.
//
// It is immutable: the only way to obtain one is via NewEmail, which trims
// surrounding whitespace, lowercases the address, and validates its format.
// Once constructed an Email is guaranteed well-formed, so callers never have to
// re-validate. The zero value is intentionally invalid (empty string) and never
// produced by NewEmail.
type Email struct {
	value string
}

// NewEmail normalizes (trim + lowercase) and validates an email address.
// It returns domain.ErrValidation if the address is empty or malformed.
func NewEmail(raw string) (Email, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	if normalized == "" {
		return Email{}, domain.ErrValidation
	}

	addr, err := mail.ParseAddress(normalized)
	if err != nil {
		return Email{}, domain.ErrValidation
	}
	// mail.ParseAddress accepts display-name forms like "Foo <a@b.com>"; we want
	// the bare address only, and it must equal what we parsed.
	if addr.Address != normalized {
		return Email{}, domain.ErrValidation
	}

	return Email{value: normalized}, nil
}

// String returns the normalized email address.
func (e Email) String() string {
	return e.value
}

// IsZero reports whether the Email is the (invalid) zero value.
func (e Email) IsZero() bool {
	return e.value == ""
}
