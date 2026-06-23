package billing

import (
	"fmt"
	"strings"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
)

// PlanInterval is a value object representing the billing cadence of a plan as a
// validated enum, replacing bare interval strings. Only month, year, and
// lifetime are representable; ParsePlanInterval rejects anything else, making
// an invalid interval unconstructable.
type PlanInterval struct {
	value string
}

// Predefined valid intervals. These wrap the same string values used at the
// repository/DTO boundary (IntervalMonth/Year/Lifetime).
var (
	Monthly  = PlanInterval{value: IntervalMonth}
	Yearly   = PlanInterval{value: IntervalYear}
	Lifetime = PlanInterval{value: IntervalLifetime}
)

// ParsePlanInterval validates and normalizes (trim + lowercase) an interval
// string into a PlanInterval. It returns domain.ErrValidation for unrecognized
// values.
func ParsePlanInterval(raw string) (PlanInterval, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case IntervalMonth:
		return Monthly, nil
	case IntervalYear:
		return Yearly, nil
	case IntervalLifetime:
		return Lifetime, nil
	default:
		return PlanInterval{}, fmt.Errorf("invalid plan interval %q: %w", raw, domain.ErrValidation)
	}
}

// String returns the canonical interval string (the value stored in the
// database and emitted in API responses).
func (i PlanInterval) String() string {
	return i.value
}

// IsLifetime reports whether the interval is a one-time (non-recurring) payment.
func (i PlanInterval) IsLifetime() bool {
	return i.value == IntervalLifetime
}

// IsRecurring reports whether the interval renews (month or year).
func (i PlanInterval) IsRecurring() bool {
	return i.value == IntervalMonth || i.value == IntervalYear
}

// IsZero reports whether the interval is the (invalid) zero value.
func (i PlanInterval) IsZero() bool {
	return i.value == ""
}
