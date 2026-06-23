package service

import "context"

type PaymentEvent struct {
	EventID        string // Stripe event ID (evt_...), used for idempotency
	Type           string // "checkout.session.completed", "invoice.paid", "customer.subscription.updated", "customer.subscription.deleted"
	CustomerID     string
	SubscriptionID string
	PriceID        string
	InvoiceID      string
	Amount         int64
	Currency       string
	InvoiceURL     string
}

// CheckoutMode determines whether the Stripe Checkout Session is for a
// recurring subscription or a one-time payment.
type CheckoutMode string

const (
	CheckoutModeSubscription CheckoutMode = "subscription"
	CheckoutModePayment      CheckoutMode = "payment"
)

// CheckoutSessionInfo holds data retrieved from a completed Stripe Checkout Session.
type CheckoutSessionInfo struct {
	CustomerID     string
	CustomerEmail  string
	SubscriptionID string
	PriceID        string
	Mode           string // "subscription" or "payment"
	Status         string // "complete", "open", "expired"
	AmountTotal    int64
}

type PaymentService interface {
	CreateCustomer(ctx context.Context, email, name string) (customerID string, err error)
	CreateCheckoutSession(ctx context.Context, customerID, priceID, successURL, cancelURL string) (sessionURL string, err error)
	CreateCheckoutSessionWithMode(ctx context.Context, customerID, priceID, successURL, cancelURL string, mode CheckoutMode) (sessionURL string, err error)
	CreateBillingPortalSession(ctx context.Context, customerID, returnURL string) (sessionURL string, err error)
	CancelSubscription(ctx context.Context, subscriptionID string) error
	HandleWebhook(payload []byte, signature string) (*PaymentEvent, error)

	// CreateGuestCheckoutSession creates a Stripe Checkout Session without an
	// existing customer (for demo/guest usage). The email is collected by Stripe.
	CreateGuestCheckoutSession(ctx context.Context, priceID, successURL, cancelURL string, mode CheckoutMode) (sessionURL string, err error)

	// RetrieveCheckoutSession fetches a completed Checkout Session's details.
	RetrieveCheckoutSession(ctx context.Context, sessionID string) (*CheckoutSessionInfo, error)
}
