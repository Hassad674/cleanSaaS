package stripe

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	stripego "github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/billingportal/session"
	checkout "github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/subscription"
	"github.com/stripe/stripe-go/v82/webhook"

	"github.com/hassad/boilerplateSaaS/backend/internal/port/service"
	"github.com/hassad/boilerplateSaaS/backend/pkg/ctxutil"
)

// Verify interface compliance at compile time.
var _ service.PaymentService = (*PaymentService)(nil)

// defaultCallTimeout is the fallback per-call ceiling used when a PaymentService
// is constructed without an explicit timeout (e.g. NewPaymentService).
const defaultCallTimeout = 15 * time.Second

type PaymentService struct {
	webhookSecret string
	callTimeout   time.Duration
}

func NewPaymentService(apiKey, webhookSecret string) *PaymentService {
	return NewPaymentServiceWithTimeout(apiKey, webhookSecret, defaultCallTimeout)
}

// NewPaymentServiceWithTimeout builds a PaymentService that bounds every Stripe
// SDK call to callTimeout (a ceiling; a nearer caller deadline still wins).
func NewPaymentServiceWithTimeout(apiKey, webhookSecret string, callTimeout time.Duration) *PaymentService {
	Init(apiKey)
	return &PaymentService{webhookSecret: webhookSecret, callTimeout: callTimeout}
}

// withTimeout derives a context bounded by the configured per-call timeout. The
// returned context is wired into the Stripe params (params.Context) so the SDK's
// HTTP request honors cancellation/deadlines.
func (s *PaymentService) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return ctxutil.WithTimeout(ctx, s.callTimeout)
}

func (s *PaymentService) CreateCustomer(ctx context.Context, email, name string) (string, error) {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	params := &stripego.CustomerParams{
		Email: stripego.String(email),
		Name:  stripego.String(name),
	}
	params.Context = ctx
	c, err := customer.New(params)
	if err != nil {
		return "", fmt.Errorf("creating stripe customer: %w", err)
	}
	return c.ID, nil
}

func (s *PaymentService) CreateCheckoutSession(ctx context.Context, customerID, priceID, successURL, cancelURL string) (string, error) {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	params := &stripego.CheckoutSessionParams{
		Customer: stripego.String(customerID),
		Mode:     stripego.String(string(stripego.CheckoutSessionModeSubscription)),
		LineItems: []*stripego.CheckoutSessionLineItemParams{
			{
				Price:    stripego.String(priceID),
				Quantity: stripego.Int64(1),
			},
		},
		SuccessURL: stripego.String(successURL),
		CancelURL:  stripego.String(cancelURL),
	}
	params.Context = ctx
	sess, err := checkout.New(params)
	if err != nil {
		return "", fmt.Errorf("creating checkout session: %w", err)
	}
	return sess.URL, nil
}

func (s *PaymentService) CreateCheckoutSessionWithMode(ctx context.Context, customerID, priceID, successURL, cancelURL string, mode service.CheckoutMode) (string, error) {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	stripeMode := string(stripego.CheckoutSessionModeSubscription)
	if mode == service.CheckoutModePayment {
		stripeMode = string(stripego.CheckoutSessionModePayment)
	}

	params := &stripego.CheckoutSessionParams{
		Customer: stripego.String(customerID),
		Mode:     stripego.String(stripeMode),
		LineItems: []*stripego.CheckoutSessionLineItemParams{
			{
				Price:    stripego.String(priceID),
				Quantity: stripego.Int64(1),
			},
		},
		SuccessURL: stripego.String(successURL),
		CancelURL:  stripego.String(cancelURL),
	}
	params.Context = ctx
	sess, err := checkout.New(params)
	if err != nil {
		return "", fmt.Errorf("creating checkout session with mode: %w", err)
	}
	return sess.URL, nil
}

func (s *PaymentService) CreateGuestCheckoutSession(ctx context.Context, priceID, successURL, cancelURL string, mode service.CheckoutMode) (string, error) {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	stripeMode := string(stripego.CheckoutSessionModeSubscription)
	if mode == service.CheckoutModePayment {
		stripeMode = string(stripego.CheckoutSessionModePayment)
	}

	params := &stripego.CheckoutSessionParams{
		Mode: stripego.String(stripeMode),
		LineItems: []*stripego.CheckoutSessionLineItemParams{
			{
				Price:    stripego.String(priceID),
				Quantity: stripego.Int64(1),
			},
		},
		SuccessURL: stripego.String(successURL),
		CancelURL:  stripego.String(cancelURL),
	}
	params.Context = ctx
	sess, err := checkout.New(params)
	if err != nil {
		return "", fmt.Errorf("creating guest checkout session: %w", err)
	}
	return sess.URL, nil
}

func (s *PaymentService) CreateBillingPortalSession(ctx context.Context, customerID, returnURL string) (string, error) {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	params := &stripego.BillingPortalSessionParams{
		Customer:  stripego.String(customerID),
		ReturnURL: stripego.String(returnURL),
	}
	params.Context = ctx
	sess, err := session.New(params)
	if err != nil {
		return "", fmt.Errorf("creating billing portal session: %w", err)
	}
	return sess.URL, nil
}

func (s *PaymentService) CancelSubscription(ctx context.Context, subscriptionID string) error {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	params := &stripego.SubscriptionParams{
		CancelAtPeriodEnd: stripego.Bool(true),
	}
	params.Context = ctx
	_, err := subscription.Update(subscriptionID, params)
	if err != nil {
		return fmt.Errorf("canceling subscription: %w", err)
	}
	return nil
}

func (s *PaymentService) RetrieveCheckoutSession(ctx context.Context, sessionID string) (*service.CheckoutSessionInfo, error) {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	params := &stripego.CheckoutSessionParams{}
	params.AddExpand("line_items")
	params.Context = ctx

	sess, err := checkout.Get(sessionID, params)
	if err != nil {
		return nil, fmt.Errorf("retrieving checkout session: %w", err)
	}

	info := &service.CheckoutSessionInfo{
		Status:      string(sess.Status),
		Mode:        string(sess.Mode),
		AmountTotal: sess.AmountTotal,
	}

	if sess.Customer != nil {
		info.CustomerID = sess.Customer.ID
	}
	if sess.CustomerEmail != "" {
		info.CustomerEmail = sess.CustomerEmail
	} else if sess.CustomerDetails != nil {
		info.CustomerEmail = sess.CustomerDetails.Email
	}
	if sess.Subscription != nil {
		info.SubscriptionID = sess.Subscription.ID
	}
	if sess.LineItems != nil && len(sess.LineItems.Data) > 0 && sess.LineItems.Data[0].Price != nil {
		info.PriceID = sess.LineItems.Data[0].Price.ID
	}

	return info, nil
}

func (s *PaymentService) HandleWebhook(payload []byte, signature string) (*service.PaymentEvent, error) {
	event, err := webhook.ConstructEvent(payload, signature, s.webhookSecret)
	if err != nil {
		return nil, fmt.Errorf("verifying webhook signature: %w", err)
	}

	pe := &service.PaymentEvent{EventID: event.ID, Type: string(event.Type)}

	switch event.Type {
	case "checkout.session.completed":
		var sess stripego.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &sess); err != nil {
			return nil, fmt.Errorf("parsing checkout session: %w", err)
		}
		pe.CustomerID = sess.Customer.ID
		if sess.Subscription != nil {
			pe.SubscriptionID = sess.Subscription.ID
		}

	case "invoice.paid":
		var inv stripego.Invoice
		if err := json.Unmarshal(event.Data.Raw, &inv); err != nil {
			return nil, fmt.Errorf("parsing invoice: %w", err)
		}
		pe.CustomerID = inv.Customer.ID
		pe.InvoiceID = inv.ID
		pe.Amount = inv.AmountPaid
		pe.Currency = string(inv.Currency)
		pe.InvoiceURL = inv.HostedInvoiceURL
		if inv.Parent != nil && inv.Parent.SubscriptionDetails != nil && inv.Parent.SubscriptionDetails.Subscription != nil {
			pe.SubscriptionID = inv.Parent.SubscriptionDetails.Subscription.ID
		}

	case "customer.subscription.updated", "customer.subscription.deleted":
		var sub stripego.Subscription
		if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
			return nil, fmt.Errorf("parsing subscription: %w", err)
		}
		pe.CustomerID = sub.Customer.ID
		pe.SubscriptionID = sub.ID
		if sub.Items != nil && len(sub.Items.Data) > 0 {
			pe.PriceID = sub.Items.Data[0].Price.ID
		}

	default:
		return nil, nil // unhandled event type
	}

	return pe, nil
}
