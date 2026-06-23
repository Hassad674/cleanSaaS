package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/notification"
	"github.com/hassad/boilerplateSaaS/backend/pkg/orgctx"
)

// NotificationRepository implements repository.NotificationRepository. It holds a
// DBTX so the same code runs on the pool or an org-scoped transaction. Every query
// filters by the active org_id (defense layer 2) on top of RLS (layer 3); inserts
// stamp the active org_id onto the row.
type NotificationRepository struct {
	db DBTX
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// newNotificationRepositoryTx binds the repository to an open transaction (org scope).
func newNotificationRepositoryTx(tx DBTX) *NotificationRepository {
	return &NotificationRepository{db: tx}
}

func (r *NotificationRepository) Create(ctx context.Context, n *notification.Notification) error {
	orgID, ok := orgctx.OrgID(ctx)
	if !ok {
		return fmt.Errorf("creating notification: %w", domain.ErrForbidden)
	}
	dataJSON, _ := json.Marshal(n.Data)
	if n.Data == nil {
		dataJSON = []byte("{}")
	}

	err := r.db.QueryRowContext(ctx,
		`INSERT INTO notifications (user_id, org_id, type, title, message, data)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, created_at`,
		n.UserID, orgID, n.Type, n.Title, n.Message, dataJSON,
	).Scan(&n.ID, &n.CreatedAt)
	if err != nil {
		return fmt.Errorf("creating notification: %w", err)
	}
	return nil
}

func (r *NotificationRepository) FindByID(ctx context.Context, id string) (*notification.Notification, error) {
	n := &notification.Notification{}
	var dataJSON []byte
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, type, title, message, read, data, created_at
		 FROM notifications WHERE id = $1 AND org_id = $2`, id, orgFilter(ctx),
	).Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.Read, &dataJSON, &n.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("finding notification: %w", err)
	}
	_ = json.Unmarshal(dataJSON, &n.Data)
	return n, nil
}

func (r *NotificationRepository) ListByUserID(ctx context.Context, userID string, unreadOnly bool, offset, limit int) ([]*notification.Notification, int, error) {
	org := orgFilter(ctx)
	countQuery := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND org_id = $2`
	listQuery := `SELECT id, user_id, type, title, message, read, data, created_at
		 FROM notifications WHERE user_id = $1 AND org_id = $2`

	if unreadOnly {
		countQuery += ` AND read = false`
		listQuery += ` AND read = false`
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, userID, org).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("counting notifications: %w", err)
	}

	listQuery += ` ORDER BY created_at DESC LIMIT $3 OFFSET $4`
	rows, err := r.db.QueryContext(ctx, listQuery, userID, org, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("listing notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*notification.Notification
	for rows.Next() {
		n := &notification.Notification{}
		var dataJSON []byte
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.Read, &dataJSON, &n.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scanning notification: %w", err)
		}
		_ = json.Unmarshal(dataJSON, &n.Data)
		notifications = append(notifications, n)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterating notifications: %w", err)
	}

	return notifications, total, nil
}

func (r *NotificationRepository) MarkRead(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `UPDATE notifications SET read = true WHERE id = $1 AND org_id = $2`, id, orgFilter(ctx))
	if err != nil {
		return fmt.Errorf("marking notification read: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *NotificationRepository) MarkAllRead(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE notifications SET read = true WHERE user_id = $1 AND org_id = $2 AND read = false`, userID, orgFilter(ctx))
	if err != nil {
		return fmt.Errorf("marking all notifications read: %w", err)
	}
	return nil
}

func (r *NotificationRepository) UnreadCount(ctx context.Context, userID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND org_id = $2 AND read = false`, userID, orgFilter(ctx),
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting unread notifications: %w", err)
	}
	return count, nil
}
