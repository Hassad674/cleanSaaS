package response

import (
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/billing"
)

type PlanResponse struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	PriceCents int      `json:"price_cents"`
	Interval   string   `json:"interval"`
	Features   []string `json:"features"`
}

func PlanFromDomain(p *billing.Plan) PlanResponse {
	features := p.Features
	if features == nil {
		features = []string{}
	}
	return PlanResponse{
		ID:         p.ID,
		Name:       p.Name,
		PriceCents: p.PriceCents,
		Interval:   p.Interval,
		Features:   features,
	}
}

type SubscriptionResponse struct {
	ID                string `json:"id"`
	PlanID            string `json:"plan_id"`
	Status            string `json:"status"`
	CurrentPeriodEnd  string `json:"current_period_end"`
	CancelAtPeriodEnd bool   `json:"cancel_at_period_end"`
}

func SubscriptionFromDomain(s *billing.Subscription) SubscriptionResponse {
	return SubscriptionResponse{
		ID:                s.ID,
		PlanID:            s.PlanID,
		Status:            string(s.Status),
		CurrentPeriodEnd:  s.CurrentPeriodEnd.Format("2006-01-02T15:04:05Z"),
		CancelAtPeriodEnd: s.CancelAtPeriodEnd,
	}
}

type InvoiceResponse struct {
	ID          string `json:"id"`
	AmountCents int    `json:"amount_cents"`
	Currency    string `json:"currency"`
	Status      string `json:"status"`
	InvoiceURL  string `json:"invoice_url"`
	CreatedAt   string `json:"created_at"`
}

func InvoiceFromDomain(i *billing.Invoice) InvoiceResponse {
	return InvoiceResponse{
		ID:          i.ID,
		AmountCents: i.AmountCents,
		Currency:    i.Currency,
		Status:      i.Status,
		InvoiceURL:  i.InvoiceURL,
		CreatedAt:   i.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
