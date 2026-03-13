package service

import "context"

type PaymentEvent struct {
	Type           string // "checkout.session.completed", "invoice.paid", "customer.subscription.updated", "customer.subscription.deleted"
	CustomerID     string
	SubscriptionID string
	PriceID        string
	InvoiceID      string
	Amount         int64
	Currency       string
	InvoiceURL     string
}

type PaymentService interface {
	CreateCustomer(ctx context.Context, email, name string) (customerID string, err error)
	CreateCheckoutSession(ctx context.Context, customerID, priceID, successURL, cancelURL string) (sessionURL string, err error)
	CreateBillingPortalSession(ctx context.Context, customerID, returnURL string) (sessionURL string, err error)
	CancelSubscription(ctx context.Context, subscriptionID string) error
	HandleWebhook(payload []byte, signature string) (*PaymentEvent, error)
}
