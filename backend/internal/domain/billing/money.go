package billing

import (
	"fmt"
	"strings"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
)

// DefaultCurrency is used when a Money is constructed without an explicit
// currency. Stripe and the seed data default to USD.
const DefaultCurrency = "usd"

// Money is a value object representing a monetary amount as an integer number of
// minor units (cents) in a given currency. It never uses floating point, so it
// is free of rounding error.
//
// It is immutable and self-validating: NewMoney rejects negative amounts and
// malformed currency codes, and normalizes the currency to a lowercase
// three-letter code. The zero value (0 of empty currency) is invalid; use
// NewMoney to construct one.
type Money struct {
	cents    int64
	currency string
}

// NewMoney constructs a Money from an amount in minor units (cents) and a
// currency code. The amount must be non-negative and the currency a 3-letter
// code (case-insensitive; normalized to lowercase). An empty currency defaults
// to DefaultCurrency. Returns domain.ErrValidation on invalid input.
func NewMoney(cents int64, currency string) (Money, error) {
	if cents < 0 {
		return Money{}, fmt.Errorf("money amount cannot be negative: %w", domain.ErrValidation)
	}

	code := strings.ToLower(strings.TrimSpace(currency))
	if code == "" {
		code = DefaultCurrency
	}
	if len(code) != 3 {
		return Money{}, fmt.Errorf("currency %q must be a 3-letter code: %w", currency, domain.ErrValidation)
	}

	return Money{cents: cents, currency: code}, nil
}

// Cents returns the amount in minor units.
func (m Money) Cents() int64 {
	return m.cents
}

// Currency returns the normalized (lowercase) 3-letter currency code.
func (m Money) Currency() string {
	return m.currency
}

// IsZero reports whether the amount is zero (e.g. a free plan).
func (m Money) IsZero() bool {
	return m.cents == 0
}

// String formats the amount as a major-unit decimal with its currency code,
// e.g. Money{1900,"usd"} -> "19.00 USD".
func (m Money) String() string {
	return fmt.Sprintf("%d.%02d %s", m.cents/100, m.cents%100, strings.ToUpper(m.currency))
}
