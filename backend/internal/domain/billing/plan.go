package billing

import "time"

type Plan struct {
	ID            string
	Name          string
	StripePriceID string
	PriceCents    int
	Interval      string // "month", "year"
	Features      []string
	IsActive      bool
	SortOrder     int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
