package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/notification"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, n *notification.Notification) error {
	dataJSON, _ := json.Marshal(n.Data)
	if n.Data == nil {
		dataJSON = []byte("{}")
	}

	err := r.db.QueryRowContext(ctx,
		`INSERT INTO notifications (user_id, type, title, message, data)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, created_at`,
		n.UserID, n.Type, n.Title, n.Message, dataJSON,
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
		 FROM notifications WHERE id = $1`, id,
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
	countQuery := `SELECT COUNT(*) FROM notifications WHERE user_id = $1`
	listQuery := `SELECT id, user_id, type, title, message, read, data, created_at
		 FROM notifications WHERE user_id = $1`

	if unreadOnly {
		countQuery += ` AND read = false`
		listQuery += ` AND read = false`
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("counting notifications: %w", err)
	}

	listQuery += ` ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, listQuery, userID, limit, offset)
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
	result, err := r.db.ExecContext(ctx, `UPDATE notifications SET read = true WHERE id = $1`, id)
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
	_, err := r.db.ExecContext(ctx, `UPDATE notifications SET read = true WHERE user_id = $1 AND read = false`, userID)
	if err != nil {
		return fmt.Errorf("marking all notifications read: %w", err)
	}
	return nil
}

func (r *NotificationRepository) UnreadCount(ctx context.Context, userID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND read = false`, userID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting unread notifications: %w", err)
	}
	return count, nil
}
