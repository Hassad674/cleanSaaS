package service

import "context"

type PaymentEvent struct {
	Type           string // "subscription.created", "invoice.paid", etc.
	CustomerID     string
	SubscriptionID string
	PriceID        string
	InvoiceID      string
	Amount         int64
	Currency       string
}

type PaymentService interface {
	CreateCustomer(ctx context.Context, email, name string) (customerID string, err error)
	CreateSubscription(ctx context.Context, customerID, priceID string) (subscriptionID string, err error)
	CancelSubscription(ctx context.Context, subscriptionID string) error
	ResumeSubscription(ctx context.Context, subscriptionID string) error
	UpdateSubscription(ctx context.Context, subscriptionID, newPriceID string) error
	CreateBillingPortalURL(ctx context.Context, customerID, returnURL string) (string, error)
	HandleWebhook(payload []byte, signature string) (*PaymentEvent, error)
}
