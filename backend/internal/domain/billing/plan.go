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

type Plan struct {
	ID            string
	Name          string
	StripePriceID string
	PriceCents    int
	Interval      string // "month", "year", "lifetime"
	Features      []string
	IsActive      bool
	SortOrder     int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// IsLifetime returns true if this plan is a one-time payment (no recurring subscription).
func (p *Plan) IsLifetime() bool {
	return p.Interval == IntervalLifetime
}
