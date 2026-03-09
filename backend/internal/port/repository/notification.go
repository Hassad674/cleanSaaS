package repository

import (
	"context"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain/notification"
)

type NotificationRepository interface {
	Create(ctx context.Context, n *notification.Notification) error
	FindByID(ctx context.Context, id string) (*notification.Notification, error)
	ListByUserID(ctx context.Context, userID string, unreadOnly bool, offset, limit int) ([]*notification.Notification, int, error)
	MarkRead(ctx context.Context, id string) error
	MarkAllRead(ctx context.Context, userID string) error
}
