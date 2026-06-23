package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	domainnotif "github.com/hassad/boilerplateSaaS/backend/internal/domain/notification"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/service"
)

type Service struct {
	notifications repository.NotificationScope
	broadcaster   service.Broadcaster // optional — nil if WebSocket is not enabled
}

// NewService wires the notification use cases. notifications is an org-scoped
// unit-of-work: each notification database operation runs inside a transaction
// bound to the caller's active organization, so RLS enforces tenant isolation on
// every query (in addition to the repository's own org_id filter).
func NewService(notifications repository.NotificationScope) *Service {
	return &Service{notifications: notifications}
}

// SetBroadcaster sets an optional real-time broadcaster (e.g., WebSocket hub).
// If set, notifications are pushed to connected clients in real time.
func (s *Service) SetBroadcaster(b service.Broadcaster) {
	s.broadcaster = b
}

func (s *Service) Send(ctx context.Context, userID, notifType, title, message string) (*domainnotif.Notification, error) {
	n := &domainnotif.Notification{
		UserID:  userID,
		Type:    notifType,
		Title:   title,
		Message: message,
	}
	err := s.notifications.WithOrgNotifications(ctx, func(notifications repository.NotificationRepository) error {
		return notifications.Create(ctx, n)
	})
	if err != nil {
		return nil, fmt.Errorf("sending notification: %w", err)
	}

	// Broadcast via WebSocket if available (fire-and-forget, never blocks notification creation).
	if s.broadcaster != nil {
		s.broadcastNotification(n)
	}

	return n, nil
}

// broadcastNotification sends a real-time WebSocket event for the notification.
func (s *Service) broadcastNotification(n *domainnotif.Notification) {
	payload := map[string]interface{}{
		"id":         n.ID,
		"type":       n.Type,
		"title":      n.Title,
		"message":    n.Message,
		"read":       n.Read,
		"created_at": n.CreatedAt,
	}

	msg := map[string]interface{}{
		"type":    "notification",
		"payload": payload,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		slog.Error("ws: failed to marshal notification", slog.String("error", err.Error()))
		return
	}

	if err := s.broadcaster.SendToUser(n.UserID, data); err != nil {
		slog.Error("ws: failed to send notification", slog.String("error", err.Error()), slog.String("user_id", n.UserID))
	}
}

func (s *Service) MarkAsRead(ctx context.Context, userID, notifID string) error {
	return s.notifications.WithOrgNotifications(ctx, func(notifications repository.NotificationRepository) error {
		n, err := notifications.FindByID(ctx, notifID)
		if err != nil {
			return err
		}
		if n.UserID != userID {
			return domain.ErrForbidden
		}
		return notifications.MarkRead(ctx, notifID)
	})
}

func (s *Service) MarkAllAsRead(ctx context.Context, userID string) error {
	return s.notifications.WithOrgNotifications(ctx, func(notifications repository.NotificationRepository) error {
		return notifications.MarkAllRead(ctx, userID)
	})
}

func (s *Service) List(ctx context.Context, userID string, unreadOnly bool, offset, limit int) ([]*domainnotif.Notification, int, error) {
	var items []*domainnotif.Notification
	var total int
	err := s.notifications.WithOrgNotifications(ctx, func(notifications repository.NotificationRepository) error {
		var e error
		items, total, e = notifications.ListByUserID(ctx, userID, unreadOnly, offset, limit)
		return e
	})
	return items, total, err
}

func (s *Service) UnreadCount(ctx context.Context, userID string) (int, error) {
	var count int
	err := s.notifications.WithOrgNotifications(ctx, func(notifications repository.NotificationRepository) error {
		var e error
		count, e = notifications.UnreadCount(ctx, userID)
		return e
	})
	return count, err
}
