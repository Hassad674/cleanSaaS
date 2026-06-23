-- Processed Stripe webhook events (idempotency ledger).
-- Records each Stripe event ID exactly once so retried deliveries are skipped.
-- No cross-feature foreign keys: this table is self-contained.
CREATE TABLE IF NOT EXISTS processed_events (
    event_id TEXT PRIMARY KEY,
    event_type TEXT NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
