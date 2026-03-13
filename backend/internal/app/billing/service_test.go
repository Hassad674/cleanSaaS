package billing

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	domainbilling "github.com/hassad/boilerplateSaaS/backend/internal/domain/billing"
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

func (m *mockUserRepo) Create(_ context.Context, _ *user.User) error           { return nil }
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

type mockPaymentSvc struct {
	createCustomerFn  func(ctx context.Context, email, name string) (string, error)
	createCheckoutFn  func(ctx context.Context, customerID, priceID, successURL, cancelURL string) (string, error)
	createPortalFn    func(ctx context.Context, customerID, returnURL string) (string, error)
	cancelSubFn       func(ctx context.Context, subscriptionID string) error
	handleWebhookFn   func(payload []byte, signature string) (*service.PaymentEvent, error)
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
			return &domainbilling.Plan{ID: "p1", StripePriceID: "price_pro"}, nil
		},
	}
	paymentSvc := &mockPaymentSvc{
		createCheckoutFn: func(_ context.Context, _, _, _, _ string) (string, error) {
			checkoutURL = "https://checkout.stripe.com/session123"
			return checkoutURL, nil
		},
	}

	svc := NewService(userRepo, &mockSubRepo{}, planRepo, &mockInvoiceRepo{}, paymentSvc, "http://localhost:3006")
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
			return &domainbilling.Plan{ID: "p1", StripePriceID: "price_pro"}, nil
		},
	}
	paymentSvc := &mockPaymentSvc{
		createCustomerFn: func(_ context.Context, _, _ string) (string, error) {
			customerCreated = true
			return "cus_new", nil
		},
	}

	svc := NewService(userRepo, &mockSubRepo{}, planRepo, &mockInvoiceRepo{}, paymentSvc, "http://localhost:3006")
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

	svc := NewService(&mockUserRepo{}, subRepo, &mockPlanRepo{}, &mockInvoiceRepo{}, paymentSvc, "http://localhost:3006")
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

	svc := NewService(&mockUserRepo{}, subRepo, &mockPlanRepo{}, &mockInvoiceRepo{}, &mockPaymentSvc{}, "http://localhost:3006")
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

	svc := NewService(&mockUserRepo{}, subRepo, &mockPlanRepo{}, &mockInvoiceRepo{}, paymentSvc, "http://localhost:3006")
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

	svc := NewService(&mockUserRepo{}, subRepo, &mockPlanRepo{}, invoiceRepo, paymentSvc, "http://localhost:3006")
	err := svc.HandleWebhook(context.Background(), []byte("payload"), "sig")
	assert.NoError(t, err)
	assert.True(t, invoiceCreated)
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

	svc := NewService(&mockUserRepo{}, subRepo, planRepo, &mockInvoiceRepo{}, paymentSvc, "http://localhost:3006")
	err := svc.HandleWebhook(context.Background(), []byte("payload"), "sig")
	assert.NoError(t, err)
	assert.Equal(t, "p_new", updatedSub.PlanID)
	assert.Equal(t, domainbilling.StatusActive, updatedSub.Status)
}

// Helper to skip the unused import warning
var _ repository.SubscriptionRepository = (*mockSubRepo)(nil)
