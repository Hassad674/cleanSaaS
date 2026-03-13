package response

import (
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/ai"
)

type ConversationResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type MessageResponse struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func ConversationFromDomain(c *ai.Conversation) ConversationResponse {
	return ConversationResponse{
		ID:        c.ID,
		Title:     c.Title,
		CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: c.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
