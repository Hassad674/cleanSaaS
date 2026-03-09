package billing

type Plan struct {
	ID            string
	Name          string
	StripePriceID string
	Price         int64  // in cents
	Currency      string // "usd", "eur"
	Interval      string // "month", "year"
	Features      []string
	Active        bool
}
