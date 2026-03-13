package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/billing"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(ctx context.Context, s *billing.Subscription) error {
	query := `INSERT INTO subscriptions (user_id, plan_id, stripe_subscription_id, status, current_period_start, current_period_end, cancel_at_period_end) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`
	err := r.db.QueryRowContext(ctx, query,
		s.UserID, s.PlanID, s.StripeSubscription, s.Status,
		s.CurrentPeriodStart, s.CurrentPeriodEnd, s.CancelAtPeriodEnd,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return fmt.Errorf("inserting subscription: %w", err)
	}
	return nil
}

func (r *SubscriptionRepository) FindByID(ctx context.Context, id string) (*billing.Subscription, error) {
	query := `SELECT id, user_id, plan_id, stripe_subscription_id, status, current_period_start, current_period_end, cancel_at_period_end, created_at, updated_at FROM subscriptions WHERE id = $1`
	return r.scanSubscription(r.db.QueryRowContext(ctx, query, id))
}

func (r *SubscriptionRepository) FindByUserID(ctx context.Context, userID string) (*billing.Subscription, error) {
	query := `SELECT id, user_id, plan_id, stripe_subscription_id, status, current_period_start, current_period_end, cancel_at_period_end, created_at, updated_at FROM subscriptions WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1`
	return r.scanSubscription(r.db.QueryRowContext(ctx, query, userID))
}

func (r *SubscriptionRepository) FindByStripeID(ctx context.Context, stripeID string) (*billing.Subscription, error) {
	query := `SELECT id, user_id, plan_id, stripe_subscription_id, status, current_period_start, current_period_end, cancel_at_period_end, created_at, updated_at FROM subscriptions WHERE stripe_subscription_id = $1`
	return r.scanSubscription(r.db.QueryRowContext(ctx, query, stripeID))
}

func (r *SubscriptionRepository) Update(ctx context.Context, s *billing.Subscription) error {
	query := `UPDATE subscriptions SET plan_id = $1, status = $2, current_period_start = $3, current_period_end = $4, cancel_at_period_end = $5, updated_at = NOW() WHERE id = $6`
	_, err := r.db.ExecContext(ctx, query, s.PlanID, s.Status, s.CurrentPeriodStart, s.CurrentPeriodEnd, s.CancelAtPeriodEnd, s.ID)
	if err != nil {
		return fmt.Errorf("updating subscription: %w", err)
	}
	return nil
}

func (r *SubscriptionRepository) scanSubscription(row *sql.Row) (*billing.Subscription, error) {
	s := &billing.Subscription{}
	err := row.Scan(&s.ID, &s.UserID, &s.PlanID, &s.StripeSubscription, &s.Status, &s.CurrentPeriodStart, &s.CurrentPeriodEnd, &s.CancelAtPeriodEnd, &s.CreatedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scanning subscription: %w", err)
	}
	return s, nil
}
