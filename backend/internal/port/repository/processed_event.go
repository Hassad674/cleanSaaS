package repository

import "context"

// ProcessedEventRepository records externally-delivered events (e.g. Stripe
// webhook events) so that retried deliveries can be detected and skipped.
type ProcessedEventRepository interface {
	// MarkProcessed atomically records that eventID has been seen. It returns
	// alreadyProcessed=true when the event had already been recorded (no new
	// row was inserted), and false when this call recorded it for the first time.
	MarkProcessed(ctx context.Context, eventID, eventType string) (alreadyProcessed bool, err error)
}
