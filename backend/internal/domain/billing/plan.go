package billing

import "time"

// Interval constants for billing plans.
const (
	IntervalMonth    = "month"
	IntervalYear     = "year"
	IntervalLifetime = "lifetime"
)

// ValidInterval returns true if the given interval is a recognized billing interval.
func ValidInterval(interval string) bool {
	switch interval {
	case IntervalMonth, IntervalYear, IntervalLifetime:
		return true
	}
	return false
}

// Plan is an intentionally anemic, validated entity (see ADR 0004): plans are
// fundamentally "validate, store, fetch" with no state machine. PriceCents and
// Interval stay as primitive columns at the repository/DTO boundary; Price() and
// IntervalVO() expose them as the Money and PlanInterval value objects when an
// invariant or formatting needs the richer type.
type Plan struct {
	ID            string
	Name          string
	StripePriceID string
	PriceCents    int
	Interval      string // canonical: see PlanInterval ("month", "year", "lifetime")
	Features      []string
	IsActive      bool
	SortOrder     int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Price returns the plan's price as a Money value object (in DefaultCurrency).
func (p *Plan) Price() (Money, error) {
	return NewMoney(int64(p.PriceCents), DefaultCurrency)
}

// IntervalVO returns the plan's interval as the PlanInterval value object,
// validating the stored string.
func (p *Plan) IntervalVO() (PlanInterval, error) {
	return ParsePlanInterval(p.Interval)
}

// IsLifetime returns true if this plan is a one-time payment (no recurring subscription).
func (p *Plan) IsLifetime() bool {
	return p.Interval == IntervalLifetime
}
