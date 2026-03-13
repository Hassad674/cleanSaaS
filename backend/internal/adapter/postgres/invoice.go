package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain/billing"
)

type InvoiceRepository struct {
	db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

func (r *InvoiceRepository) Create(ctx context.Context, i *billing.Invoice) error {
	query := `INSERT INTO invoices (user_id, stripe_invoice_id, amount_cents, currency, status, invoice_url) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`
	err := r.db.QueryRowContext(ctx, query, i.UserID, i.StripeInvoiceID, i.AmountCents, i.Currency, i.Status, i.InvoiceURL).Scan(&i.ID, &i.CreatedAt)
	if err != nil {
		return fmt.Errorf("inserting invoice: %w", err)
	}
	return nil
}

func (r *InvoiceRepository) ListByUserID(ctx context.Context, userID string, offset, limit int) ([]*billing.Invoice, int, error) {
	countQuery := `SELECT COUNT(*) FROM invoices WHERE user_id = $1`
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("counting invoices: %w", err)
	}

	query := `SELECT id, user_id, stripe_invoice_id, amount_cents, currency, status, invoice_url, created_at FROM invoices WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("listing invoices: %w", err)
	}
	defer rows.Close()

	var invoices []*billing.Invoice
	for rows.Next() {
		i := &billing.Invoice{}
		if err := rows.Scan(&i.ID, &i.UserID, &i.StripeInvoiceID, &i.AmountCents, &i.Currency, &i.Status, &i.InvoiceURL, &i.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scanning invoice: %w", err)
		}
		invoices = append(invoices, i)
	}
	return invoices, total, rows.Err()
}
