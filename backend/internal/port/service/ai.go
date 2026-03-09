package service

import (
	"context"
	"io"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain/ai"
)

type AIResponse struct {
	Content string
	Model   ai.Model
	Tokens  int
}

type AIService interface {
	Chat(ctx context.Context, messages []ai.Message) (*AIResponse, error)
	Stream(ctx context.Context, messages []ai.Message, writer io.Writer) error
}
