package billing

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	domainbilling "github.com/hassad/boilerplateSaaS/backend/internal/domain/billing"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/service"
)

type Service struct {
	users           repository.UserRepository
	subscriptions   repository.SubscriptionRepository
	plans           repository.PlanRepository
	invoices        repository.InvoiceRepository
	processedEvents repository.ProcessedEventRepository
	payment         service.PaymentService
	frontendURL     string
}

func NewService(
	users repository.UserRepository,
	subscriptions repository.SubscriptionRepository,
	plans repository.PlanRepository,
	invoices repository.InvoiceRepository,
	processedEvents repository.ProcessedEventRepository,
	payment service.PaymentService,
	frontendURL string,
) *Service {
	return &Service{
		users:           users,
		subscriptions:   subscriptions,
		plans:           plans,
		invoices:        invoices,
		processedEvents: processedEvents,
		payment:         payment,
		frontendURL:     frontendURL,
	}
}

func (s *Service) GetPlans(ctx context.Context) ([]*domainbilling.Plan, error) {
	return s.plans.List(ctx)
}

func (s *Service) CreateCheckout(ctx context.Context, userID, planID string) (string, error) {
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("finding user: %w", err)
	}

	plan, err := s.plans.FindByID(ctx, planID)
	if err != nil {
		return "", fmt.Errorf("finding plan: %w", err)
	}

	customerID := u.StripeID
	if customerID == "" {
		customerID, err = s.payment.CreateCustomer(ctx, u.Email, u.Name)
		if err != nil {
			return "", fmt.Errorf("creating stripe customer: %w", err)
		}
		u.StripeID = customerID
		u.UpdatedAt = time.Now()
		if err := s.users.Update(ctx, u); err != nil {
			return "", fmt.Errorf("saving stripe customer ID: %w", err)
		}
	}

	successURL := fmt.Sprintf("%s/settings/billing?success=true", s.frontendURL)
	cancelURL := fmt.Sprintf("%s/pricing?canceled=true", s.frontendURL)

	// Lifetime plans are one-time payments, not subscriptions
	mode := service.CheckoutModeSubscription
	if plan.IsLifetime() {
		mode = service.CheckoutModePayment
	}

	return s.payment.CreateCheckoutSessionWithMode(ctx, customerID, plan.StripePriceID, successURL, cancelURL, mode)
}

// DemoCheckout creates a Stripe Checkout Session for the public demo page.
// No user account is required — Stripe will collect the email on its checkout page.
func (s *Service) DemoCheckout(ctx context.Context, planID, successURL, cancelURL string) (string, error) {
	plan, err := s.plans.FindByID(ctx, planID)
	if err != nil {
		return "", fmt.Errorf("finding plan: %w", err)
	}

	if plan.PriceCents == 0 {
		return "", domain.ErrValidation
	}

	mode := service.CheckoutModeSubscription
	if plan.IsLifetime() {
		mode = service.CheckoutModePayment
	}

	return s.payment.CreateGuestCheckoutSession(ctx, plan.StripePriceID, successURL, cancelURL, mode)
}

// DemoSessionInfo holds the data returned by GetDemoSession for the frontend demo.
type DemoSessionInfo struct {
	PlanName      string `json:"plan_name"`
	PriceCents    int    `json:"price_cents"`
	Interval      string `json:"interval"`
	CustomerID    string `json:"customer_id"`
	CustomerEmail string `json:"customer_email"`
	Mode          string `json:"mode"`
	Status        string `json:"status"`
}

// GetDemoSession retrieves a completed Stripe Checkout Session and enriches it
// with plan details from the database.
func (s *Service) GetDemoSession(ctx context.Context, sessionID string) (*DemoSessionInfo, error) {
	sess, err := s.payment.RetrieveCheckoutSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("retrieving session: %w", err)
	}

	info := &DemoSessionInfo{
		CustomerID:    sess.CustomerID,
		CustomerEmail: sess.CustomerEmail,
		Mode:          sess.Mode,
		Status:        sess.Status,
	}

	if sess.PriceID != "" {
		plan, err := s.plans.FindByStripePriceID(ctx, sess.PriceID)
		if err == nil {
			info.PlanName = plan.Name
			info.PriceCents = plan.PriceCents
			info.Interval = plan.Interval
		}
	}

	return info, nil
}

// DemoPortalSession creates a Stripe Billing Portal session for a demo customer.
func (s *Service) DemoPortalSession(ctx context.Context, customerID, returnURL string) (string, error) {
	return s.payment.CreateBillingPortalSession(ctx, customerID, returnURL)
}

func (s *Service) GetSubscription(ctx context.Context, userID string) (*domainbilling.Subscription, error) {
	return s.subscriptions.FindByUserID(ctx, userID)
}

func (s *Service) CancelSubscription(ctx context.Context, userID string) error {
	sub, err := s.subscriptions.FindByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("finding subscription: %w", err)
	}

	if !sub.CanCancel() {
		return domain.ErrValidation
	}

	if err := s.payment.CancelSubscription(ctx, sub.StripeSubscription); err != nil {
		return fmt.Errorf("canceling stripe subscription: %w", err)
	}

	sub.Cancel()
	return s.subscriptions.Update(ctx, sub)
}

func (s *Service) CreatePortalSession(ctx context.Context, userID string) (string, error) {
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("finding user: %w", err)
	}

	if u.StripeID == "" {
		return "", domain.ErrNotFound
	}

	returnURL := fmt.Sprintf("%s/settings/billing", s.frontendURL)
	return s.payment.CreateBillingPortalSession(ctx, u.StripeID, returnURL)
}

func (s *Service) GetInvoices(ctx context.Context, userID string, offset, limit int) ([]*domainbilling.Invoice, int, error) {
	return s.invoices.ListByUserID(ctx, userID, offset, limit)
}

func (s *Service) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	event, err := s.payment.HandleWebhook(payload, signature)
	if err != nil {
		return fmt.Errorf("parsing webhook: %w", err)
	}

	if event == nil {
		return nil
	}

	// Idempotency: skip events Stripe has already delivered. We record the
	// event ID before processing so a retried delivery is a no-op. An empty
	// EventID (e.g. in tests or non-Stripe sources) falls through to normal
	// processing.
	if event.EventID != "" {
		alreadyProcessed, err := s.processedEvents.MarkProcessed(ctx, event.EventID, event.Type)
		if err != nil {
			return fmt.Errorf("recording processed event: %w", err)
		}
		if alreadyProcessed {
			slog.Info("skipping already-processed stripe webhook event",
				slog.String("event_id", event.EventID),
				slog.String("event_type", event.Type),
			)
			return nil
		}
	}

	switch event.Type {
	case "invoice.paid":
		return s.handleInvoicePaid(ctx, event)
	case "customer.subscription.updated":
		return s.handleSubscriptionUpdated(ctx, event)
	case "customer.subscription.deleted":
		return s.handleSubscriptionDeleted(ctx, event)
	}

	return nil
}

func (s *Service) handleInvoicePaid(ctx context.Context, event *service.PaymentEvent) error {
	if event.SubscriptionID == "" {
		return nil
	}

	sub, err := s.subscriptions.FindByStripeID(ctx, event.SubscriptionID)
	if err != nil {
		return nil
	}

	invoice := &domainbilling.Invoice{
		UserID:          sub.UserID,
		StripeInvoiceID: event.InvoiceID,
		AmountCents:     int(event.Amount),
		Currency:        event.Currency,
		Status:          "paid",
		InvoiceURL:      event.InvoiceURL,
	}

	return s.invoices.Create(ctx, invoice)
}

func (s *Service) handleSubscriptionUpdated(ctx context.Context, event *service.PaymentEvent) error {
	existing, err := s.subscriptions.FindByStripeID(ctx, event.SubscriptionID)
	if err != nil {
		return s.createSubscriptionFromWebhook(ctx, event)
	}

	if event.PriceID != "" {
		plan, err := s.plans.FindByStripePriceID(ctx, event.PriceID)
		if err == nil {
			existing.PlanID = plan.ID
		}
	}
	existing.Status = domainbilling.StatusActive
	existing.UpdatedAt = time.Now()

	return s.subscriptions.Update(ctx, existing)
}

func (s *Service) handleSubscriptionDeleted(ctx context.Context, event *service.PaymentEvent) error {
	sub, err := s.subscriptions.FindByStripeID(ctx, event.SubscriptionID)
	if err != nil {
		return nil
	}

	sub.Status = domainbilling.StatusCanceled
	sub.UpdatedAt = time.Now()
	return s.subscriptions.Update(ctx, sub)
}

func (s *Service) createSubscriptionFromWebhook(ctx context.Context, event *service.PaymentEvent) error {
	if event.PriceID == "" || event.CustomerID == "" {
		return nil
	}

	plan, err := s.plans.FindByStripePriceID(ctx, event.PriceID)
	if err != nil {
		return nil
	}

	// Look up user by Stripe customer ID (indexed query, not full table scan)
	u, err := s.users.FindByStripeID(ctx, event.CustomerID)
	if err != nil {
		return nil
	}

	sub := &domainbilling.Subscription{
		UserID:             u.ID,
		PlanID:             plan.ID,
		StripeSubscription: event.SubscriptionID,
		Status:             domainbilling.StatusActive,
		CurrentPeriodStart: time.Now(),
		CurrentPeriodEnd:   time.Now().Add(30 * 24 * time.Hour),
	}

	return s.subscriptions.Create(ctx, sub)
}
