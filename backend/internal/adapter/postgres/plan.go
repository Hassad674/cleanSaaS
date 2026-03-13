package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/billing"
)

type PlanRepository struct {
	db *sql.DB
}

func NewPlanRepository(db *sql.DB) *PlanRepository {
	return &PlanRepository{db: db}
}

func (r *PlanRepository) FindByID(ctx context.Context, id string) (*billing.Plan, error) {
	query := `SELECT id, name, stripe_price_id, price_cents, interval, features, is_active, sort_order, created_at, updated_at FROM plans WHERE id = $1`
	return r.scanPlan(r.db.QueryRowContext(ctx, query, id))
}

func (r *PlanRepository) FindByStripePriceID(ctx context.Context, priceID string) (*billing.Plan, error) {
	query := `SELECT id, name, stripe_price_id, price_cents, interval, features, is_active, sort_order, created_at, updated_at FROM plans WHERE stripe_price_id = $1`
	return r.scanPlan(r.db.QueryRowContext(ctx, query, priceID))
}

func (r *PlanRepository) List(ctx context.Context) ([]*billing.Plan, error) {
	query := `SELECT id, name, stripe_price_id, price_cents, interval, features, is_active, sort_order, created_at, updated_at FROM plans WHERE is_active = true ORDER BY sort_order ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("listing plans: %w", err)
	}
	defer rows.Close()

	var plans []*billing.Plan
	for rows.Next() {
		p, err := r.scanPlanRow(rows)
		if err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	return plans, rows.Err()
}

func (r *PlanRepository) scanPlan(row *sql.Row) (*billing.Plan, error) {
	p := &billing.Plan{}
	var featuresJSON []byte
	err := row.Scan(&p.ID, &p.Name, &p.StripePriceID, &p.PriceCents, &p.Interval, &featuresJSON, &p.IsActive, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scanning plan: %w", err)
	}
	if err := json.Unmarshal(featuresJSON, &p.Features); err != nil {
		return nil, fmt.Errorf("unmarshalling plan features: %w", err)
	}
	return p, nil
}

func (r *PlanRepository) scanPlanRow(rows *sql.Rows) (*billing.Plan, error) {
	p := &billing.Plan{}
	var featuresJSON []byte
	err := rows.Scan(&p.ID, &p.Name, &p.StripePriceID, &p.PriceCents, &p.Interval, &featuresJSON, &p.IsActive, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("scanning plan row: %w", err)
	}
	if err := json.Unmarshal(featuresJSON, &p.Features); err != nil {
		return nil, fmt.Errorf("unmarshalling plan features: %w", err)
	}
	return p, nil
}
