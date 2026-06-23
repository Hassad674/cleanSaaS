package billing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSubscription_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		status   Status
		expected bool
	}{
		{"Active subscription", StatusActive, true},
		{"Trialing subscription", StatusTrialing, true},
		{"Canceled subscription", StatusCanceled, false},
		{"Past due subscription", StatusPastDue, false},
		{"Inactive subscription", StatusInactive, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := &Subscription{Status: tt.status}
			assert.Equal(t, tt.expected, sub.IsActive())
		})
	}
}

func TestSubscription_CanCancel(t *testing.T) {
	tests := []struct {
		name              string
		status            Status
		cancelAtPeriodEnd bool
		expected          bool
	}{
		{"Active, not canceling", StatusActive, false, true},
		{"Trialing, not canceling", StatusTrialing, false, true},
		{"Active, already canceling", StatusActive, true, false},
		{"Canceled, not canceling", StatusCanceled, false, false},
		{"Past due, not canceling", StatusPastDue, false, false},
		{"Inactive, not canceling", StatusInactive, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := &Subscription{
				Status:            tt.status,
				CancelAtPeriodEnd: tt.cancelAtPeriodEnd,
			}
			assert.Equal(t, tt.expected, sub.CanCancel())
		})
	}
}

func TestSubscription_Cancel(t *testing.T) {
	sub := &Subscription{
		ID:                "s1",
		Status:            StatusActive,
		CancelAtPeriodEnd: false,
		UpdatedAt:         time.Now().Add(-time.Hour),
	}

	before := sub.UpdatedAt
	err := sub.Cancel()

	assert.NoError(t, err)
	assert.True(t, sub.CancelAtPeriodEnd, "should set CancelAtPeriodEnd to true")
	assert.True(t, sub.UpdatedAt.After(before), "should update UpdatedAt")
}

func TestStatus_Constants(t *testing.T) {
	assert.Equal(t, Status("active"), StatusActive)
	assert.Equal(t, Status("canceled"), StatusCanceled)
	assert.Equal(t, Status("past_due"), StatusPastDue)
	assert.Equal(t, Status("trialing"), StatusTrialing)
	assert.Equal(t, Status("inactive"), StatusInactive)
}

func TestPlan_Fields(t *testing.T) {
	plan := Plan{
		ID:            "p1",
		Name:          "Pro",
		StripePriceID: "price_pro",
		PriceCents:    1900,
		Interval:      "month",
		Features:      []string{"feature1", "feature2"},
		IsActive:      true,
		SortOrder:     1,
	}

	assert.Equal(t, "Pro", plan.Name)
	assert.Equal(t, 1900, plan.PriceCents)
	assert.Equal(t, "month", plan.Interval)
	assert.Len(t, plan.Features, 2)
	assert.True(t, plan.IsActive)
}

func TestPlan_IsLifetime(t *testing.T) {
	tests := []struct {
		name     string
		interval string
		expected bool
	}{
		{"Monthly plan", IntervalMonth, false},
		{"Yearly plan", IntervalYear, false},
		{"Lifetime plan", IntervalLifetime, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := &Plan{Interval: tt.interval}
			assert.Equal(t, tt.expected, plan.IsLifetime())
		})
	}
}

func TestValidInterval(t *testing.T) {
	assert.True(t, ValidInterval("month"))
	assert.True(t, ValidInterval("year"))
	assert.True(t, ValidInterval("lifetime"))
	assert.False(t, ValidInterval("weekly"))
	assert.False(t, ValidInterval(""))
}

func TestInvoice_Fields(t *testing.T) {
	invoice := Invoice{
		ID:              "inv-1",
		UserID:          "u1",
		StripeInvoiceID: "inv_stripe",
		AmountCents:     1900,
		Currency:        "usd",
		Status:          "paid",
		InvoiceURL:      "https://stripe.com/inv",
	}

	assert.Equal(t, "u1", invoice.UserID)
	assert.Equal(t, 1900, invoice.AmountCents)
	assert.Equal(t, "usd", invoice.Currency)
	assert.Equal(t, "paid", invoice.Status)
}
