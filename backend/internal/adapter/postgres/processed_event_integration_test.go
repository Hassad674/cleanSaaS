//go:build integration

package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProcessedEventRepository_MarkProcessed_Idempotency proves the webhook
// idempotency ledger: the first MarkProcessed for an event records it and reports
// alreadyProcessed=false; a second call for the same event id inserts nothing
// (ON CONFLICT DO NOTHING, RowsAffected==0) and reports alreadyProcessed=true.
// The unique event id is deleted afterwards so the ledger is left untouched.
func TestProcessedEventRepository_MarkProcessed_Idempotency(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	repo := NewProcessedEventRepository(db)

	eventID := "evt_" + uniqueTag()
	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(), `DELETE FROM processed_events WHERE event_id = $1`, eventID)
	})

	// First delivery: newly recorded.
	already, err := repo.MarkProcessed(ctx, eventID, "invoice.paid")
	require.NoError(t, err)
	assert.False(t, already, "first MarkProcessed records the event (not already processed)")

	// Retried delivery (same id): detected as already processed, no new row.
	already, err = repo.MarkProcessed(ctx, eventID, "invoice.paid")
	require.NoError(t, err)
	assert.True(t, already, "second MarkProcessed for the same event reports already-processed")

	// A different event id is independent and recorded fresh.
	otherID := "evt_" + uniqueTag()
	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(), `DELETE FROM processed_events WHERE event_id = $1`, otherID)
	})
	already, err = repo.MarkProcessed(ctx, otherID, "customer.subscription.updated")
	require.NoError(t, err)
	assert.False(t, already, "a distinct event id is processed for the first time")
}
