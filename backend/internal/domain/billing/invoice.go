package billing

import "time"

type Invoice struct {
	ID              string
	UserID          string
	StripeInvoiceID string
	Amount          int64
	Currency        string
	Status          string // "paid", "open", "void"
	PaidAt          *time.Time
	CreatedAt       time.Time
}
