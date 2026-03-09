package repository

import (
	"context"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain/ai"
)

type ConversationRepository interface {
	Create(ctx context.Context, c *ai.Conversation) error
	FindByID(ctx context.Context, id string) (*ai.Conversation, error)
	Update(ctx context.Context, c *ai.Conversation) error
	Delete(ctx context.Context, id string) error
	ListByUserID(ctx context.Context, userID string, offset, limit int) ([]*ai.Conversation, int, error)
}
