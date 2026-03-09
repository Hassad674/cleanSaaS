package repository

import (
	"context"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain/billing"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, s *billing.Subscription) error
	FindByID(ctx context.Context, id string) (*billing.Subscription, error)
	FindByUserID(ctx context.Context, userID string) (*billing.Subscription, error)
	Update(ctx context.Context, s *billing.Subscription) error
}

type PlanRepository interface {
	FindByID(ctx context.Context, id string) (*billing.Plan, error)
	FindByStripePriceID(ctx context.Context, priceID string) (*billing.Plan, error)
	List(ctx context.Context) ([]*billing.Plan, error)
}

type InvoiceRepository interface {
	Create(ctx context.Context, i *billing.Invoice) error
	ListByUserID(ctx context.Context, userID string, offset, limit int) ([]*billing.Invoice, int, error)
}
