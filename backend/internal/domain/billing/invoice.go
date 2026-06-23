package billing

import "time"

// Invoice is an intentionally anemic, validated record (see ADR 0004): it is an
// immutable receipt with no state machine. AmountCents and Currency stay as
// primitive columns at the repository/DTO boundary; Amount() exposes them as the
// Money value object.
type Invoice struct {
	ID              string
	UserID          string
	StripeInvoiceID string
	AmountCents     int
	Currency        string
	Status          string // "paid", "open", "void"
	InvoiceURL      string
	CreatedAt       time.Time
}

// Amount returns the invoice total as a Money value object.
func (i *Invoice) Amount() (Money, error) {
	return NewMoney(int64(i.AmountCents), i.Currency)
}
