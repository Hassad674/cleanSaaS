package notification

import (
	"context"
	"fmt"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	domainnotif "github.com/hassad/boilerplateSaaS/backend/internal/domain/notification"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
)

type Service struct {
	notifications repository.NotificationRepository
}

func NewService(notifications repository.NotificationRepository) *Service {
	return &Service{notifications: notifications}
}

func (s *Service) Send(ctx context.Context, userID, notifType, title, message string) (*domainnotif.Notification, error) {
	n := &domainnotif.Notification{
		UserID:  userID,
		Type:    notifType,
		Title:   title,
		Message: message,
	}
	if err := s.notifications.Create(ctx, n); err != nil {
		return nil, fmt.Errorf("sending notification: %w", err)
	}
	return n, nil
}

func (s *Service) MarkAsRead(ctx context.Context, userID, notifID string) error {
	n, err := s.notifications.FindByID(ctx, notifID)
	if err != nil {
		return err
	}
	if n.UserID != userID {
		return domain.ErrForbidden
	}
	return s.notifications.MarkRead(ctx, notifID)
}

func (s *Service) MarkAllAsRead(ctx context.Context, userID string) error {
	return s.notifications.MarkAllRead(ctx, userID)
}

func (s *Service) List(ctx context.Context, userID string, unreadOnly bool, offset, limit int) ([]*domainnotif.Notification, int, error) {
	return s.notifications.ListByUserID(ctx, userID, unreadOnly, offset, limit)
}

func (s *Service) UnreadCount(ctx context.Context, userID string) (int, error) {
	return s.notifications.UnreadCount(ctx, userID)
}
