package stripe

import (
	"context"
	"encoding/json"
	"fmt"

	stripego "github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/billingportal/session"
	checkout "github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/subscription"
	"github.com/stripe/stripe-go/v82/webhook"

	"github.com/hassad/boilerplateSaaS/backend/internal/port/service"
)

// Verify interface compliance at compile time.
var _ service.PaymentService = (*PaymentService)(nil)

type PaymentService struct {
	webhookSecret string
}

func NewPaymentService(apiKey, webhookSecret string) *PaymentService {
	Init(apiKey)
	return &PaymentService{webhookSecret: webhookSecret}
}

func (s *PaymentService) CreateCustomer(_ context.Context, email, name string) (string, error) {
	params := &stripego.CustomerParams{
		Email: stripego.String(email),
		Name:  stripego.String(name),
	}
	c, err := customer.New(params)
	if err != nil {
		return "", fmt.Errorf("creating stripe customer: %w", err)
	}
	return c.ID, nil
}

func (s *PaymentService) CreateCheckoutSession(_ context.Context, customerID, priceID, successURL, cancelURL string) (string, error) {
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
	sess, err := checkout.New(params)
	if err != nil {
		return "", fmt.Errorf("creating checkout session: %w", err)
	}
	return sess.URL, nil
}

func (s *PaymentService) CreateCheckoutSessionWithMode(_ context.Context, customerID, priceID, successURL, cancelURL string, mode service.CheckoutMode) (string, error) {
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
	sess, err := checkout.New(params)
	if err != nil {
		return "", fmt.Errorf("creating checkout session with mode: %w", err)
	}
	return sess.URL, nil
}

func (s *PaymentService) CreateGuestCheckoutSession(_ context.Context, priceID, successURL, cancelURL string, mode service.CheckoutMode) (string, error) {
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
	sess, err := checkout.New(params)
	if err != nil {
		return "", fmt.Errorf("creating guest checkout session: %w", err)
	}
	return sess.URL, nil
}

func (s *PaymentService) CreateBillingPortalSession(_ context.Context, customerID, returnURL string) (string, error) {
	params := &stripego.BillingPortalSessionParams{
		Customer:  stripego.String(customerID),
		ReturnURL: stripego.String(returnURL),
	}
	sess, err := session.New(params)
	if err != nil {
		return "", fmt.Errorf("creating billing portal session: %w", err)
	}
	return sess.URL, nil
}

func (s *PaymentService) CancelSubscription(_ context.Context, subscriptionID string) error {
	params := &stripego.SubscriptionParams{
		CancelAtPeriodEnd: stripego.Bool(true),
	}
	_, err := subscription.Update(subscriptionID, params)
	if err != nil {
		return fmt.Errorf("canceling subscription: %w", err)
	}
	return nil
}

func (s *PaymentService) RetrieveCheckoutSession(_ context.Context, sessionID string) (*service.CheckoutSessionInfo, error) {
	params := &stripego.CheckoutSessionParams{}
	params.AddExpand("line_items")

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

	pe := &service.PaymentEvent{Type: string(event.Type)}

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
