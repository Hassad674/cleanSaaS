package billing

import "time"

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
