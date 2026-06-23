package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

type ProcessedEventRepository struct {
	db *sql.DB
}

func NewProcessedEventRepository(db *sql.DB) *ProcessedEventRepository {
	return &ProcessedEventRepository{db: db}
}

// MarkProcessed records eventID via INSERT ... ON CONFLICT DO NOTHING. When the
// row already exists no insert happens (RowsAffected == 0), so the event was
// already processed by a prior delivery and the caller should skip it.
func (r *ProcessedEventRepository) MarkProcessed(ctx context.Context, eventID, eventType string) (bool, error) {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO processed_events (event_id, event_type)
		 VALUES ($1, $2)
		 ON CONFLICT (event_id) DO NOTHING`,
		eventID, eventType,
	)
	if err != nil {
		return false, fmt.Errorf("marking event as processed: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("checking rows affected: %w", err)
	}

	alreadyProcessed := rows == 0
	return alreadyProcessed, nil
}
