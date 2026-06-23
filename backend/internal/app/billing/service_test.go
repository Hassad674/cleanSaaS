package billing

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	domainbilling "github.com/hassad/boilerplateSaaS/backend/internal/domain/billing"
	domainorg "github.com/hassad/boilerplateSaaS/backend/internal/domain/org"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/user"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/service"
)

// Mocks

type mockUserRepo struct {
	findByIDFn func(ctx context.Context, id string) (*user.User, error)
	updateFn   func(ctx context.Context, u *user.User) error
	listFn     func(ctx context.Context, offset, limit int) ([]*user.User, int, error)
}

func (m *mockUserRepo) Create(_ context.Context, _ *user.User) error { return nil }
func (m *mockUserRepo) FindByEmail(_ context.Context, _ string) (*user.User, error) {
	return nil, domain.ErrNotFound
}
func (m *mockUserRepo) FindByID(ctx context.Context, id string) (*user.User, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, domain.ErrNotFound
}
func (m *mockUserRepo) FindByProvider(_ context.Context, _, _ string) (*user.User, error) {
	return nil, domain.ErrNotFound
}
func (m *mockUserRepo) Update(ctx context.Context, u *user.User) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, u)
	}
	return nil
}
func (m *mockUserRepo) FindByStripeID(_ context.Context, _ string) (*user.User, error) {
	return nil, domain.ErrNotFound
}
func (m *mockUserRepo) Delete(_ context.Context, _ string) error { return nil }
func (m *mockUserRepo) List(ctx context.Context, offset, limit int) ([]*user.User, int, error) {
	if m.listFn != nil {
		return m.listFn(ctx, offset, limit)
	}
	return nil, 0, nil
}
func (m *mockUserRepo) Search(_ context.Context, _ string, _, _ int) ([]*user.User, int, error) {
	return nil, 0, nil
}
func (m *mockUserRepo) Count(_ context.Context) (int, error) { return 0, nil }

type mockSubRepo struct {
	createFn       func(ctx context.Context, s *domainbilling.Subscription) error
	findByUserIDFn func(ctx context.Context, userID string) (*domainbilling.Subscription, error)
	findByStripeFn func(ctx context.Context, stripeID string) (*domainbilling.Subscription, error)
	updateFn       func(ctx context.Context, s *domainbilling.Subscription) error
}

func (m *mockSubRepo) Create(ctx context.Context, s *domainbilling.Subscription) error {
	if m.createFn != nil {
		return m.createFn(ctx, s)
	}
	return nil
}
func (m *mockSubRepo) FindByID(_ context.Context, _ string) (*domainbilling.Subscription, error) {
	return nil, domain.ErrNotFound
}
func (m *mockSubRepo) FindByUserID(ctx context.Context, userID string) (*domainbilling.Subscription, error) {
	if m.findByUserIDFn != nil {
		return m.findByUserIDFn(ctx, userID)
	}
	return nil, domain.ErrNotFound
}
func (m *mockSubRepo) FindByStripeID(ctx context.Context, stripeID string) (*domainbilling.Subscription, error) {
	if m.findByStripeFn != nil {
		return m.findByStripeFn(ctx, stripeID)
	}
	return nil, domain.ErrNotFound
}
func (m *mockSubRepo) Update(ctx context.Context, s *domainbilling.Subscription) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, s)
	}
	return nil
}

type mockPlanRepo struct {
	findByIDFn    func(ctx context.Context, id string) (*domainbilling.Plan, error)
	findByPriceFn func(ctx context.Context, priceID string) (*domainbilling.Plan, error)
	listFn        func(ctx context.Context) ([]*domainbilling.Plan, error)
}

func (m *mockPlanRepo) FindByID(ctx context.Context, id string) (*domainbilling.Plan, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, domain.ErrNotFound
}
func (m *mockPlanRepo) FindByStripePriceID(ctx context.Context, priceID string) (*domainbilling.Plan, error) {
	if m.findByPriceFn != nil {
		return m.findByPriceFn(ctx, priceID)
	}
	return nil, domain.ErrNotFound
}
func (m *mockPlanRepo) List(ctx context.Context) ([]*domainbilling.Plan, error) {
	if m.listFn != nil {
		return m.listFn(ctx)
	}
	return nil, nil
}

type mockInvoiceRepo struct {
	createFn func(ctx context.Context, i *domainbilling.Invoice) error
}

func (m *mockInvoiceRepo) Create(ctx context.Context, i *domainbilling.Invoice) error {
	if m.createFn != nil {
		return m.createFn(ctx, i)
	}
	return nil
}
func (m *mockInvoiceRepo) ListByUserID(_ context.Context, _ string, _, _ int) ([]*domainbilling.Invoice, int, error) {
	return nil, 0, nil
}

type mockProcessedEventRepo struct {
	markProcessedFn func(ctx context.Context, eventID, eventType string) (bool, error)
}

func (m *mockProcessedEventRepo) MarkProcessed(ctx context.Context, eventID, eventType string) (bool, error) {
	if m.markProcessedFn != nil {
		return m.markProcessedFn(ctx, eventID, eventType)
	}
	return false, nil
}

type mockPaymentSvc struct {
	createCustomerFn      func(ctx context.Context, email, name string) (string, error)
	createCheckoutFn      func(ctx context.Context, customerID, priceID, successURL, cancelURL string) (string, error)
	createCheckoutModeFn  func(ctx context.Context, customerID, priceID, successURL, cancelURL string, mode service.CheckoutMode) (string, error)
	createGuestCheckoutFn func(ctx context.Context, priceID, successURL, cancelURL string, mode service.CheckoutMode) (string, error)
	createPortalFn        func(ctx context.Context, customerID, returnURL string) (string, error)
	cancelSubFn           func(ctx context.Context, subscriptionID string) error
	handleWebhookFn       func(payload []byte, signature string) (*service.PaymentEvent, error)
	retrieveSessionFn     func(ctx context.Context, sessionID string) (*service.CheckoutSessionInfo, error)
}

func (m *mockPaymentSvc) CreateCustomer(ctx context.Context, email, name string) (string, error) {
	if m.createCustomerFn != nil {
		return m.createCustomerFn(ctx, email, name)
	}
	return "cus_test", nil
}
func (m *mockPaymentSvc) CreateCheckoutSession(ctx context.Context, customerID, priceID, successURL, cancelURL string) (string, error) {
	if m.createCheckoutFn != nil {
		return m.createCheckoutFn(ctx, customerID, priceID, successURL, cancelURL)
	}
	return "https://checkout.stripe.com/test", nil
}
func (m *mockPaymentSvc) CreateCheckoutSessionWithMode(ctx context.Context, customerID, priceID, successURL, cancelURL string, mode service.CheckoutMode) (string, error) {
	if m.createCheckoutModeFn != nil {
		return m.createCheckoutModeFn(ctx, customerID, priceID, successURL, cancelURL, mode)
	}
	return "https://checkout.stripe.com/test", nil
}
func (m *mockPaymentSvc) CreateGuestCheckoutSession(ctx context.Context, priceID, successURL, cancelURL string, mode service.CheckoutMode) (string, error) {
	if m.createGuestCheckoutFn != nil {
		return m.createGuestCheckoutFn(ctx, priceID, successURL, cancelURL, mode)
	}
	return "https://checkout.stripe.com/guest-test", nil
}
func (m *mockPaymentSvc) CreateBillingPortalSession(ctx context.Context, customerID, returnURL string) (string, error) {
	if m.createPortalFn != nil {
		return m.createPortalFn(ctx, customerID, returnURL)
	}
	return "https://billing.stripe.com/test", nil
}
func (m *mockPaymentSvc) CancelSubscription(ctx context.Context, subscriptionID string) error {
	if m.cancelSubFn != nil {
		return m.cancelSubFn(ctx, subscriptionID)
	}
	return nil
}
func (m *mockPaymentSvc) HandleWebhook(payload []byte, signature string) (*service.PaymentEvent, error) {
	if m.handleWebhookFn != nil {
		return m.handleWebhookFn(payload, signature)
	}
	return nil, nil
}
func (m *mockPaymentSvc) RetrieveCheckoutSession(ctx context.Context, sessionID string) (*service.CheckoutSessionInfo, error) {
	if m.retrieveSessionFn != nil {
		return m.retrieveSessionFn(ctx, sessionID)
	}
	return &service.CheckoutSessionInfo{}, nil
}

// Tests

func TestBillingService_CreateCheckout_Success(t *testing.T) {
	var checkoutURL string
	userRepo := &mockUserRepo{
		findByIDFn: func(_ context.Context, _ string) (*user.User, error) {
			return &user.User{ID: "u1", Email: "test@test.com", Name: "Test", StripeID: "cus_123"}, nil
		},
	}
	planRepo := &mockPlanRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainbilling.Plan, error) {
			return &domainbilling.Plan{ID: "p1", StripePriceID: "price_pro", Interval: "month"}, nil
		},
	}
	paymentSvc := &mockPaymentSvc{
		createCheckoutModeFn: func(_ context.Context, _, _, _, _ string, mode service.CheckoutMode) (string, error) {
			assert.Equal(t, service.CheckoutModeSubscription, mode)
			checkoutURL = "https://checkout.stripe.com/session123"
			return checkoutURL, nil
		},
	}

	svc := newTestService(userRepo, &mockSubRepo{}, planRepo, &mockInvoiceRepo{}, &mockProcessedEventRepo{}, paymentSvc, "http://localhost:3006")
	url, err := svc.CreateCheckout(context.Background(), "u1", "p1")
	assert.NoError(t, err)
	assert.Equal(t, "https://checkout.stripe.com/session123", url)
}

func TestBillingService_CreateCheckout_CreatesStripeCustomer(t *testing.T) {
	var customerCreated bool
	var userUpdated bool
	userRepo := &mockUserRepo{
		findByIDFn: func(_ context.Context, _ string) (*user.User, error) {
			return &user.User{ID: "u1", Email: "new@test.com", Name: "New", StripeID: ""}, nil
		},
		updateFn: func(_ context.Context, u *user.User) error {
			userUpdated = true
			assert.Equal(t, "cus_new", u.StripeID)
			return nil
		},
	}
	planRepo := &mockPlanRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainbilling.Plan, error) {
			return &domainbilling.Plan{ID: "p1", StripePriceID: "price_pro", Interval: "month"}, nil
		},
	}
	paymentSvc := &mockPaymentSvc{
		createCustomerFn: func(_ context.Context, _, _ string) (string, error) {
			customerCreated = true
			return "cus_new", nil
		},
	}

	svc := newTestService(userRepo, &mockSubRepo{}, planRepo, &mockInvoiceRepo{}, &mockProcessedEventRepo{}, paymentSvc, "http://localhost:3006")
	_, err := svc.CreateCheckout(context.Background(), "u1", "p1")
	assert.NoError(t, err)
	assert.True(t, customerCreated)
	assert.True(t, userUpdated)
}

func TestBillingService_CancelSubscription_Success(t *testing.T) {
	var stripeCanceled bool
	var subUpdated bool
	subRepo := &mockSubRepo{
		findByUserIDFn: func(_ context.Context, _ string) (*domainbilling.Subscription, error) {
			return &domainbilling.Subscription{
				ID:                 "s1",
				UserID:             "u1",
				StripeSubscription: "sub_123",
				Status:             domainbilling.StatusActive,
			}, nil
		},
		updateFn: func(_ context.Context, s *domainbilling.Subscription) error {
			subUpdated = true
			assert.True(t, s.CancelAtPeriodEnd)
			return nil
		},
	}
	paymentSvc := &mockPaymentSvc{
		cancelSubFn: func(_ context.Context, _ string) error {
			stripeCanceled = true
			return nil
		},
	}

	svc := newTestService(&mockUserRepo{}, subRepo, &mockPlanRepo{}, &mockInvoiceRepo{}, &mockProcessedEventRepo{}, paymentSvc, "http://localhost:3006")
	err := svc.CancelSubscription(context.Background(), "u1")
	assert.NoError(t, err)
	assert.True(t, stripeCanceled)
	assert.True(t, subUpdated)
}

func TestBillingService_CancelSubscription_AlreadyCanceling(t *testing.T) {
	subRepo := &mockSubRepo{
		findByUserIDFn: func(_ context.Context, _ string) (*domainbilling.Subscription, error) {
			return &domainbilling.Subscription{
				ID:                "s1",
				Status:            domainbilling.StatusActive,
				CancelAtPeriodEnd: true,
			}, nil
		},
	}

	svc := newTestService(&mockUserRepo{}, subRepo, &mockPlanRepo{}, &mockInvoiceRepo{}, &mockProcessedEventRepo{}, &mockPaymentSvc{}, "http://localhost:3006")
	err := svc.CancelSubscription(context.Background(), "u1")
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestBillingService_HandleWebhook_SubscriptionDeleted(t *testing.T) {
	var subStatus domainbilling.Status
	subRepo := &mockSubRepo{
		findByStripeFn: func(_ context.Context, _ string) (*domainbilling.Subscription, error) {
			return &domainbilling.Subscription{
				ID:                 "s1",
				StripeSubscription: "sub_123",
				Status:             domainbilling.StatusActive,
			}, nil
		},
		updateFn: func(_ context.Context, s *domainbilling.Subscription) error {
			subStatus = s.Status
			return nil
		},
	}
	paymentSvc := &mockPaymentSvc{
		handleWebhookFn: func(_ []byte, _ string) (*service.PaymentEvent, error) {
			return &service.PaymentEvent{
				Type:           "customer.subscription.deleted",
				SubscriptionID: "sub_123",
				CustomerID:     "cus_123",
			}, nil
		},
	}

	svc := newTestService(&mockUserRepo{}, subRepo, &mockPlanRepo{}, &mockInvoiceRepo{}, &mockProcessedEventRepo{}, paymentSvc, "http://localhost:3006")
	err := svc.HandleWebhook(context.Background(), []byte("payload"), "sig")
	assert.NoError(t, err)
	assert.Equal(t, domainbilling.StatusCanceled, subStatus)
}

func TestBillingService_HandleWebhook_InvoicePaid(t *testing.T) {
	var invoiceCreated bool
	subRepo := &mockSubRepo{
		findByStripeFn: func(_ context.Context, _ string) (*domainbilling.Subscription, error) {
			return &domainbilling.Subscription{ID: "s1", UserID: "u1"}, nil
		},
	}
	invoiceRepo := &mockInvoiceRepo{
		createFn: func(_ context.Context, i *domainbilling.Invoice) error {
			invoiceCreated = true
			assert.Equal(t, "u1", i.UserID)
			assert.Equal(t, int(1900), i.AmountCents)
			assert.Equal(t, "paid", i.Status)
			return nil
		},
	}
	paymentSvc := &mockPaymentSvc{
		handleWebhookFn: func(_ []byte, _ string) (*service.PaymentEvent, error) {
			return &service.PaymentEvent{
				Type:           "invoice.paid",
				SubscriptionID: "sub_123",
				InvoiceID:      "inv_123",
				Amount:         1900,
				Currency:       "usd",
			}, nil
		},
	}

	svc := newTestService(&mockUserRepo{}, subRepo, &mockPlanRepo{}, invoiceRepo, &mockProcessedEventRepo{}, paymentSvc, "http://localhost:3006")
	err := svc.HandleWebhook(context.Background(), []byte("payload"), "sig")
	assert.NoError(t, err)
	assert.True(t, invoiceCreated)
}

// TestBillingService_HandleWebhook_Idempotent proves that the same Stripe event
// delivered twice (Stripe's automatic retry) is processed exactly once: the
// second delivery is recognised as already-processed and short-circuits before
// any business logic runs.
func TestBillingService_HandleWebhook_Idempotent(t *testing.T) {
	var invoiceCreateCount int
	subRepo := &mockSubRepo{
		findByStripeFn: func(_ context.Context, _ string) (*domainbilling.Subscription, error) {
			return &domainbilling.Subscription{ID: "s1", UserID: "u1"}, nil
		},
	}
	invoiceRepo := &mockInvoiceRepo{
		createFn: func(_ context.Context, _ *domainbilling.Invoice) error {
			invoiceCreateCount++
			return nil
		},
	}
	paymentSvc := &mockPaymentSvc{
		handleWebhookFn: func(_ []byte, _ string) (*service.PaymentEvent, error) {
			return &service.PaymentEvent{
				EventID:        "evt_123",
				Type:           "invoice.paid",
				SubscriptionID: "sub_123",
				InvoiceID:      "inv_123",
				Amount:         1900,
				Currency:       "usd",
			}, nil
		},
	}

	// Stateful processed-event store mirroring INSERT ... ON CONFLICT DO NOTHING:
	// the first MarkProcessed inserts (alreadyProcessed=false), subsequent calls
	// for the same event report alreadyProcessed=true.
	seen := map[string]bool{}
	processedRepo := &mockProcessedEventRepo{
		markProcessedFn: func(_ context.Context, eventID, _ string) (bool, error) {
			if seen[eventID] {
				return true, nil
			}
			seen[eventID] = true
			return false, nil
		},
	}

	svc := newTestService(&mockUserRepo{}, subRepo, &mockPlanRepo{}, invoiceRepo, processedRepo, paymentSvc, "http://localhost:3006")

	// First delivery: processed, invoice created.
	err := svc.HandleWebhook(context.Background(), []byte("payload"), "sig")
	assert.NoError(t, err)

	// Second delivery of the SAME event (Stripe retry): no-op.
	err = svc.HandleWebhook(context.Background(), []byte("payload"), "sig")
	assert.NoError(t, err)

	assert.Equal(t, 1, invoiceCreateCount, "invoice must be created only once across duplicate deliveries")
}

func TestBillingService_HandleWebhook_SubscriptionUpdated_Existing(t *testing.T) {
	var updatedSub *domainbilling.Subscription
	subRepo := &mockSubRepo{
		findByStripeFn: func(_ context.Context, _ string) (*domainbilling.Subscription, error) {
			return &domainbilling.Subscription{
				ID:                 "s1",
				UserID:             "u1",
				PlanID:             "p_old",
				StripeSubscription: "sub_123",
				Status:             domainbilling.StatusActive,
				CurrentPeriodEnd:   time.Now().Add(30 * 24 * time.Hour),
			}, nil
		},
		updateFn: func(_ context.Context, s *domainbilling.Subscription) error {
			updatedSub = s
			return nil
		},
	}
	planRepo := &mockPlanRepo{
		findByPriceFn: func(_ context.Context, _ string) (*domainbilling.Plan, error) {
			return &domainbilling.Plan{ID: "p_new"}, nil
		},
	}
	paymentSvc := &mockPaymentSvc{
		handleWebhookFn: func(_ []byte, _ string) (*service.PaymentEvent, error) {
			return &service.PaymentEvent{
				Type:           "customer.subscription.updated",
				SubscriptionID: "sub_123",
				PriceID:        "price_new",
				CustomerID:     "cus_123",
			}, nil
		},
	}

	svc := newTestService(&mockUserRepo{}, subRepo, planRepo, &mockInvoiceRepo{}, &mockProcessedEventRepo{}, paymentSvc, "http://localhost:3006")
	err := svc.HandleWebhook(context.Background(), []byte("payload"), "sig")
	assert.NoError(t, err)
	assert.Equal(t, "p_new", updatedSub.PlanID)
	assert.Equal(t, domainbilling.StatusActive, updatedSub.Status)
}

func TestBillingService_CreateCheckout_LifetimePlan(t *testing.T) {
	var receivedMode service.CheckoutMode
	userRepo := &mockUserRepo{
		findByIDFn: func(_ context.Context, _ string) (*user.User, error) {
			return &user.User{ID: "u1", Email: "test@test.com", Name: "Test", StripeID: "cus_123"}, nil
		},
	}
	planRepo := &mockPlanRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainbilling.Plan, error) {
			return &domainbilling.Plan{ID: "p1", StripePriceID: "price_lifetime", Interval: "lifetime"}, nil
		},
	}
	paymentSvc := &mockPaymentSvc{
		createCheckoutModeFn: func(_ context.Context, _, _, _, _ string, mode service.CheckoutMode) (string, error) {
			receivedMode = mode
			return "https://checkout.stripe.com/lifetime", nil
		},
	}

	svc := newTestService(userRepo, &mockSubRepo{}, planRepo, &mockInvoiceRepo{}, &mockProcessedEventRepo{}, paymentSvc, "http://localhost:3006")
	url, err := svc.CreateCheckout(context.Background(), "u1", "p1")
	assert.NoError(t, err)
	assert.Equal(t, "https://checkout.stripe.com/lifetime", url)
	assert.Equal(t, service.CheckoutModePayment, receivedMode)
}

func TestBillingService_DemoCheckout_Success(t *testing.T) {
	var receivedMode service.CheckoutMode
	planRepo := &mockPlanRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainbilling.Plan, error) {
			return &domainbilling.Plan{ID: "p1", StripePriceID: "price_pro_monthly", PriceCents: 1900, Interval: "month"}, nil
		},
	}
	paymentSvc := &mockPaymentSvc{
		createGuestCheckoutFn: func(_ context.Context, _, _, _ string, mode service.CheckoutMode) (string, error) {
			receivedMode = mode
			return "https://checkout.stripe.com/guest", nil
		},
	}

	svc := newTestService(&mockUserRepo{}, &mockSubRepo{}, planRepo, &mockInvoiceRepo{}, &mockProcessedEventRepo{}, paymentSvc, "http://localhost:3006")
	url, err := svc.DemoCheckout(context.Background(), "p1", "http://localhost/success", "http://localhost/cancel")
	assert.NoError(t, err)
	assert.Equal(t, "https://checkout.stripe.com/guest", url)
	assert.Equal(t, service.CheckoutModeSubscription, receivedMode)
}

func TestBillingService_DemoCheckout_FreePlanRejected(t *testing.T) {
	planRepo := &mockPlanRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainbilling.Plan, error) {
			return &domainbilling.Plan{ID: "free", StripePriceID: "price_free", PriceCents: 0, Interval: "month"}, nil
		},
	}

	svc := newTestService(&mockUserRepo{}, &mockSubRepo{}, planRepo, &mockInvoiceRepo{}, &mockProcessedEventRepo{}, &mockPaymentSvc{}, "http://localhost:3006")
	_, err := svc.DemoCheckout(context.Background(), "free", "http://localhost/success", "http://localhost/cancel")
	assert.Error(t, err)
}

func TestBillingService_DemoCheckout_LifetimeMode(t *testing.T) {
	var receivedMode service.CheckoutMode
	planRepo := &mockPlanRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainbilling.Plan, error) {
			return &domainbilling.Plan{ID: "p1", StripePriceID: "price_lifetime", PriceCents: 49900, Interval: "lifetime"}, nil
		},
	}
	paymentSvc := &mockPaymentSvc{
		createGuestCheckoutFn: func(_ context.Context, _, _, _ string, mode service.CheckoutMode) (string, error) {
			receivedMode = mode
			return "https://checkout.stripe.com/guest-lifetime", nil
		},
	}

	svc := newTestService(&mockUserRepo{}, &mockSubRepo{}, planRepo, &mockInvoiceRepo{}, &mockProcessedEventRepo{}, paymentSvc, "http://localhost:3006")
	url, err := svc.DemoCheckout(context.Background(), "p1", "http://localhost/success", "http://localhost/cancel")
	assert.NoError(t, err)
	assert.Equal(t, "https://checkout.stripe.com/guest-lifetime", url)
	assert.Equal(t, service.CheckoutModePayment, receivedMode)
}

// mockSubScope adapts a plain subscription repo into a repository.SubscriptionScope
// by invoking the callback with the underlying repo — unit tests need no real tx.
type mockSubScope struct {
	repo repository.SubscriptionRepository
}

func (s *mockSubScope) WithOrgSubscriptions(ctx context.Context, fn func(subscriptions repository.SubscriptionRepository) error) error {
	return fn(s.repo)
}

// mockOrgRepo is a minimal OrganizationRepository for billing tests. The webhook
// path resolves a customer's org via FindDefaultForUser; it returns a fixed org.
type mockOrgRepo struct {
	findDefaultFn func(ctx context.Context, userID string) (*domainorg.Organization, error)
}

func (m *mockOrgRepo) Create(_ context.Context, _ *domainorg.Organization) error { return nil }
func (m *mockOrgRepo) FindByID(_ context.Context, _ string) (*domainorg.Organization, error) {
	return nil, domain.ErrNotFound
}
func (m *mockOrgRepo) FindBySlug(_ context.Context, _ string) (*domainorg.Organization, error) {
	return nil, domain.ErrNotFound
}
func (m *mockOrgRepo) FindDefaultForUser(ctx context.Context, userID string) (*domainorg.Organization, error) {
	if m.findDefaultFn != nil {
		return m.findDefaultFn(ctx, userID)
	}
	return &domainorg.Organization{ID: "org-1", OwnerID: userID}, nil
}

// newTestService builds the billing service from the positional dependencies the
// tests already construct, wiring the subscription repo into BOTH the raw seam
// (Subscriptions, used by the webhook/system path) and the org-scoped seam
// (SubscriptionTx, used by the authenticated request path).
func newTestService(
	users repository.UserRepository,
	subs repository.SubscriptionRepository,
	plans repository.PlanRepository,
	invoices repository.InvoiceRepository,
	processed repository.ProcessedEventRepository,
	payment service.PaymentService,
	frontendURL string,
) *Service {
	return NewService(Deps{
		Users:           users,
		Orgs:            &mockOrgRepo{},
		Subscriptions:   subs,
		SubscriptionTx:  &mockSubScope{repo: subs},
		Plans:           plans,
		Invoices:        invoices,
		ProcessedEvents: processed,
		Payment:         payment,
		FrontendURL:     frontendURL,
	})
}

// Compile-time checks that the mocks satisfy their port interfaces.
var (
	_ repository.SubscriptionRepository   = (*mockSubRepo)(nil)
	_ repository.ProcessedEventRepository = (*mockProcessedEventRepo)(nil)
	_ repository.SubscriptionScope        = (*mockSubScope)(nil)
	_ repository.OrganizationRepository   = (*mockOrgRepo)(nil)
)
