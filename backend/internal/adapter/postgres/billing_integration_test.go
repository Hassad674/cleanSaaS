//go:build integration

package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/billing"
)

// firstSeededPlan returns a seeded, active plan from the catalog (read-only). The
// suite never mutates plans; subscriptions just reference one by id.
func firstSeededPlan(ctx context.Context, t *testing.T, repo *PlanRepository) *billing.Plan {
	t.Helper()
	plans, err := repo.List(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, plans, "the DB must have seeded plans (run `make seed`)")
	return plans[0]
}

// TestPlanRepository_SeededReads proves the read-only plan catalog: List returns
// active plans ordered by sort_order with their JSONB features unmarshalled, and
// FindByID / FindByStripePriceID resolve the same row. It never writes.
func TestPlanRepository_SeededReads(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	repo := NewPlanRepository(db)

	plans, err := repo.List(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, plans, "seeded plans expected")

	// List is ordered by sort_order ascending.
	for i := 1; i < len(plans); i++ {
		assert.LessOrEqual(t, plans[i-1].SortOrder, plans[i].SortOrder, "plans ordered by sort_order")
	}

	p := plans[0]
	assert.True(t, p.IsActive)
	assert.NotNil(t, p.Features, "JSONB features column unmarshals (never nil)")

	byID, err := repo.FindByID(ctx, p.ID)
	require.NoError(t, err)
	assert.Equal(t, p.Name, byID.Name)

	if p.StripePriceID != "" {
		byPrice, err := repo.FindByStripePriceID(ctx, p.StripePriceID)
		require.NoError(t, err)
		assert.Equal(t, p.ID, byPrice.ID)
	}
}

// TestPlanRepository_NotFound proves missing-plan lookups map to not-found.
func TestPlanRepository_NotFound(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	repo := NewPlanRepository(db)

	_, err := repo.FindByID(ctx, uniqueTag())
	assert.ErrorIs(t, err, domain.ErrNotFound)
	_, err = repo.FindByStripePriceID(ctx, "price_missing_"+uniqueTag())
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

// TestSubscriptionRepository_SystemPathRoundTrip exercises the subscription
// write path as the SYSTEM (webhook) path does: on the privileged pool with
// org_id supplied explicitly on the aggregate. It covers Create (RETURNING),
// FindByID, FindByStripeID, and Update. The org-scoped READ path is covered in
// the tenant-scope integration test.
func TestSubscriptionRepository_SystemPathRoundTrip(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	repo := NewSubscriptionRepository(db)
	planRepo := NewPlanRepository(db)

	owner := newUser(ctx, t, db)
	o := newOrg(ctx, t, db, owner.ID)
	plan := firstSeededPlan(ctx, t, planRepo)

	now := time.Now()
	s := &billing.Subscription{
		UserID:             owner.ID,
		OrgID:              o.ID,
		PlanID:             plan.ID,
		StripeSubscription: "sub_" + uniqueTag(),
		Status:             billing.StatusActive,
		CurrentPeriodStart: now,
		CurrentPeriodEnd:   now.Add(30 * 24 * time.Hour),
		CancelAtPeriodEnd:  false,
	}
	require.NoError(t, repo.Create(ctx, s))
	require.NotEmpty(t, s.ID, "Create populates id via RETURNING")
	t.Cleanup(func() { _, _ = db.ExecContext(context.Background(), `DELETE FROM subscriptions WHERE id = $1`, s.ID) })

	byID, err := repo.FindByID(ctx, s.ID)
	require.NoError(t, err)
	assert.Equal(t, o.ID, byID.OrgID)
	assert.Equal(t, billing.StatusActive, byID.Status)

	byStripe, err := repo.FindByStripeID(ctx, s.StripeSubscription)
	require.NoError(t, err)
	assert.Equal(t, s.ID, byStripe.ID)

	// Update mutates status + cancel flag.
	s.Status = billing.StatusCanceled
	s.CancelAtPeriodEnd = true
	require.NoError(t, repo.Update(ctx, s))

	updated, err := repo.FindByID(ctx, s.ID)
	require.NoError(t, err)
	assert.Equal(t, billing.StatusCanceled, updated.Status)
	assert.True(t, updated.CancelAtPeriodEnd)
}

// TestSubscriptionRepository_Create_MissingOrg proves the guard that rejects a
// subscription without an org_id with a validation error (no row written).
func TestSubscriptionRepository_Create_MissingOrg(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	repo := NewSubscriptionRepository(db)

	err := repo.Create(ctx, &billing.Subscription{
		StripeSubscription: "sub_noorg_" + uniqueTag(),
		Status:             billing.StatusActive,
	})
	assert.ErrorIs(t, err, domain.ErrValidation, "subscription without org is rejected")
}

// TestInvoiceRepository_CreateAndListPagination exercises invoice Create
// (RETURNING id, created_at) and the ListByUserID COUNT + LIMIT/OFFSET path
// scoped to a single user, asserting both the page slice and the total.
func TestInvoiceRepository_CreateAndListPagination(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	repo := NewInvoiceRepository(db)

	owner := newUser(ctx, t, db)

	const want = 3
	for i := 0; i < want; i++ {
		inv := &billing.Invoice{
			UserID:          owner.ID,
			StripeInvoiceID: "in_" + uniqueTag(),
			AmountCents:     1000 + i,
			Currency:        "usd",
			Status:          "paid",
			InvoiceURL:      "https://example.test/" + uniqueTag(),
		}
		require.NoError(t, repo.Create(ctx, inv))
		require.NotEmpty(t, inv.ID)
		require.False(t, inv.CreatedAt.IsZero())
	}

	// Page 1: 2 of the 3, with the correct total. Deleting the user cascades to
	// invoices, so no explicit invoice cleanup is needed.
	page, total, err := repo.ListByUserID(ctx, owner.ID, 0, 2)
	require.NoError(t, err)
	assert.Equal(t, want, total, "COUNT(*) scoped to the user equals the 3 created")
	assert.Len(t, page, 2, "LIMIT 2 yields a 2-row page")

	// Page 2: the remaining 1.
	page2, _, err := repo.ListByUserID(ctx, owner.ID, 2, 2)
	require.NoError(t, err)
	assert.Len(t, page2, 1, "OFFSET 2 yields the final row")
}
