package ai

import (
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/service"
)

type Service struct {
	conversations repository.ConversationRepository
	ai            service.AIService
}

func NewService(conversations repository.ConversationRepository, ai service.AIService) *Service {
	return &Service{conversations: conversations, ai: ai}
}

// Chat, Stream, GetHistory will be implemented
