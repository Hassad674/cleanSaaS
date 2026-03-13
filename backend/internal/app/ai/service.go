package ai

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	domainai "github.com/hassad/boilerplateSaaS/backend/internal/domain/ai"
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

func (s *Service) CreateConversation(ctx context.Context, userID, title string) (*domainai.Conversation, error) {
	if title == "" {
		title = "New conversation"
	}
	c := &domainai.Conversation{
		UserID: userID,
		Title:  title,
	}
	if err := s.conversations.Create(ctx, c); err != nil {
		return nil, fmt.Errorf("creating conversation: %w", err)
	}
	return c, nil
}

func (s *Service) SendMessage(ctx context.Context, userID, conversationID, content string) (string, error) {
	conv, err := s.conversations.FindByID(ctx, conversationID)
	if err != nil {
		return "", err
	}
	if conv.UserID != userID {
		return "", domain.ErrForbidden
	}

	// Save user message
	userMsg := domainai.Message{Role: domainai.RoleUser, Content: content}
	if err := s.conversations.AddMessage(ctx, conversationID, userMsg); err != nil {
		return "", fmt.Errorf("saving user message: %w", err)
	}

	// Build message history for AI
	conv.Messages = append(conv.Messages, userMsg)

	// Call AI
	resp, err := s.ai.Chat(ctx, conv.Messages)
	if err != nil {
		return "", fmt.Errorf("AI chat: %w", err)
	}

	// Save assistant response
	assistantMsg := domainai.Message{Role: domainai.RoleAssistant, Content: resp.Content}
	if err := s.conversations.AddMessage(ctx, conversationID, assistantMsg); err != nil {
		return "", fmt.Errorf("saving assistant message: %w", err)
	}

	// Auto-title from first message
	if len(conv.Messages) == 1 {
		title := content
		if len(title) > 50 {
			title = title[:50] + "..."
		}
		conv.Title = title
		_ = s.conversations.Update(ctx, conv)
	}

	return resp.Content, nil
}

func (s *Service) StreamMessage(ctx context.Context, userID, conversationID, content string, writer io.Writer) error {
	conv, err := s.conversations.FindByID(ctx, conversationID)
	if err != nil {
		return err
	}
	if conv.UserID != userID {
		return domain.ErrForbidden
	}

	// Save user message
	userMsg := domainai.Message{Role: domainai.RoleUser, Content: content}
	if err := s.conversations.AddMessage(ctx, conversationID, userMsg); err != nil {
		return fmt.Errorf("saving user message: %w", err)
	}

	conv.Messages = append(conv.Messages, userMsg)

	// Stream from AI, collecting full response
	collector := &responseCollector{writer: writer}
	if err := s.ai.Stream(ctx, conv.Messages, collector); err != nil {
		return fmt.Errorf("AI stream: %w", err)
	}

	// Save full assistant response
	assistantMsg := domainai.Message{Role: domainai.RoleAssistant, Content: collector.Content()}
	if err := s.conversations.AddMessage(ctx, conversationID, assistantMsg); err != nil {
		return fmt.Errorf("saving assistant message: %w", err)
	}

	// Auto-title
	if len(conv.Messages) == 1 {
		title := content
		if len(title) > 50 {
			title = title[:50] + "..."
		}
		conv.Title = title
		_ = s.conversations.Update(ctx, conv)
	}

	return nil
}

func (s *Service) GetConversation(ctx context.Context, userID, conversationID string) (*domainai.Conversation, error) {
	conv, err := s.conversations.FindByID(ctx, conversationID)
	if err != nil {
		return nil, err
	}
	if conv.UserID != userID {
		return nil, domain.ErrForbidden
	}
	return conv, nil
}

func (s *Service) ListConversations(ctx context.Context, userID string, offset, limit int) ([]*domainai.Conversation, int, error) {
	return s.conversations.ListByUserID(ctx, userID, offset, limit)
}

func (s *Service) DeleteConversation(ctx context.Context, userID, conversationID string) error {
	conv, err := s.conversations.FindByID(ctx, conversationID)
	if err != nil {
		return err
	}
	if conv.UserID != userID {
		return domain.ErrForbidden
	}
	return s.conversations.Delete(ctx, conversationID)
}

// responseCollector captures streamed content while forwarding to the writer.
// Uses strings.Builder to avoid O(n^2) string concatenation on each chunk.
type responseCollector struct {
	writer io.Writer
	buf    strings.Builder
}

func (c *responseCollector) Write(p []byte) (n int, err error) {
	c.buf.Write(p)
	return c.writer.Write(p)
}

// Content returns the accumulated streamed content.
func (c *responseCollector) Content() string {
	return c.buf.String()
}
