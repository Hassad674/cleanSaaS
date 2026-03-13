package response

import (
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/notification"
)

type NotificationResponse struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	Read      bool   `json:"read"`
	CreatedAt string `json:"created_at"`
}

func NotificationFromDomain(n *notification.Notification) NotificationResponse {
	return NotificationResponse{
		ID:        n.ID,
		Type:      n.Type,
		Title:     n.Title,
		Message:   n.Message,
		Read:      n.Read,
		CreatedAt: n.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
