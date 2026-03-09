package billing

import (
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/service"
)

type Service struct {
	users         repository.UserRepository
	subscriptions repository.SubscriptionRepository
	plans         repository.PlanRepository
	invoices      repository.InvoiceRepository
	payment       service.PaymentService
}

func NewService(
	users repository.UserRepository,
	subscriptions repository.SubscriptionRepository,
	plans repository.PlanRepository,
	invoices repository.InvoiceRepository,
	payment service.PaymentService,
) *Service {
	return &Service{
		users:         users,
		subscriptions: subscriptions,
		plans:         plans,
		invoices:      invoices,
		payment:       payment,
	}
}

// Subscribe, Cancel, ChangePlan, HandleWebhook, GetInvoices
// will be implemented when Stripe is integrated
