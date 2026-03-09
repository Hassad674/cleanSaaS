package notification

import (
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/service"
)

type Service struct {
	notifications repository.NotificationRepository
	email         service.EmailService
}

func NewService(notifications repository.NotificationRepository, email service.EmailService) *Service {
	return &Service{notifications: notifications, email: email}
}

// Send, MarkRead, GetUnread will be implemented
