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
	conversations repository.ConversationScope
	ai            service.AIService
}

// NewService wires the AI chat use cases. conversations is an org-scoped
// unit-of-work: each conversation/message database operation runs inside a
// transaction bound to the caller's active organization, so RLS enforces tenant
// isolation on every query. The external AI call is made OUTSIDE any scope so no
// database transaction is held open across a slow network request.
func NewService(conversations repository.ConversationScope, ai service.AIService) *Service {
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
	err := s.conversations.WithOrgConversations(ctx, func(conversations repository.ConversationRepository) error {
		return conversations.Create(ctx, c)
	})
	if err != nil {
		return nil, fmt.Errorf("creating conversation: %w", err)
	}
	return c, nil
}

func (s *Service) SendMessage(ctx context.Context, userID, conversationID, content string) (string, error) {
	conv, err := s.loadOwnedConversation(ctx, userID, conversationID)
	if err != nil {
		return "", err
	}

	// Save user message (org-scoped) and build history.
	userMsg := domainai.Message{Role: domainai.RoleUser, Content: content}
	if err := s.addMessage(ctx, conversationID, userMsg); err != nil {
		return "", fmt.Errorf("saving user message: %w", err)
	}
	conv.Messages = append(conv.Messages, userMsg)

	// Call AI OUTSIDE any DB transaction.
	resp, err := s.ai.Chat(ctx, conv.Messages)
	if err != nil {
		return "", fmt.Errorf("AI chat: %w", err)
	}

	assistantMsg := domainai.Message{Role: domainai.RoleAssistant, Content: resp.Content}
	if err := s.addMessage(ctx, conversationID, assistantMsg); err != nil {
		return "", fmt.Errorf("saving assistant message: %w", err)
	}

	s.autoTitle(ctx, conv, content)
	return resp.Content, nil
}

func (s *Service) StreamMessage(ctx context.Context, userID, conversationID, content string, writer io.Writer) error {
	conv, err := s.loadOwnedConversation(ctx, userID, conversationID)
	if err != nil {
		return err
	}

	userMsg := domainai.Message{Role: domainai.RoleUser, Content: content}
	if err := s.addMessage(ctx, conversationID, userMsg); err != nil {
		return fmt.Errorf("saving user message: %w", err)
	}
	conv.Messages = append(conv.Messages, userMsg)

	// Stream from AI OUTSIDE any DB transaction, collecting the full response.
	collector := &responseCollector{writer: writer}
	if err := s.ai.Stream(ctx, conv.Messages, collector); err != nil {
		return fmt.Errorf("AI stream: %w", err)
	}

	assistantMsg := domainai.Message{Role: domainai.RoleAssistant, Content: collector.Content()}
	if err := s.addMessage(ctx, conversationID, assistantMsg); err != nil {
		return fmt.Errorf("saving assistant message: %w", err)
	}

	s.autoTitle(ctx, conv, content)
	return nil
}

func (s *Service) GetConversation(ctx context.Context, userID, conversationID string) (*domainai.Conversation, error) {
	return s.loadOwnedConversation(ctx, userID, conversationID)
}

func (s *Service) ListConversations(ctx context.Context, userID string, offset, limit int) ([]*domainai.Conversation, int, error) {
	var convos []*domainai.Conversation
	var total int
	err := s.conversations.WithOrgConversations(ctx, func(repo repository.ConversationRepository) error {
		var e error
		convos, total, e = repo.ListByUserID(ctx, userID, offset, limit)
		return e
	})
	return convos, total, err
}

func (s *Service) DeleteConversation(ctx context.Context, userID, conversationID string) error {
	return s.conversations.WithOrgConversations(ctx, func(conversations repository.ConversationRepository) error {
		conv, err := conversations.FindByID(ctx, conversationID)
		if err != nil {
			return err
		}
		if conv.UserID != userID {
			return domain.ErrForbidden
		}
		return conversations.Delete(ctx, conversationID)
	})
}

// loadOwnedConversation fetches a conversation under the active org scope and
// verifies the caller owns it, returning domain.ErrForbidden otherwise.
func (s *Service) loadOwnedConversation(ctx context.Context, userID, conversationID string) (*domainai.Conversation, error) {
	var conv *domainai.Conversation
	err := s.conversations.WithOrgConversations(ctx, func(conversations repository.ConversationRepository) error {
		c, err := conversations.FindByID(ctx, conversationID)
		if err != nil {
			return err
		}
		if c.UserID != userID {
			return domain.ErrForbidden
		}
		conv = c
		return nil
	})
	if err != nil {
		return nil, err
	}
	return conv, nil
}

// addMessage persists a single message under the active org scope.
func (s *Service) addMessage(ctx context.Context, conversationID string, msg domainai.Message) error {
	return s.conversations.WithOrgConversations(ctx, func(conversations repository.ConversationRepository) error {
		return conversations.AddMessage(ctx, conversationID, msg)
	})
}

// autoTitle sets the conversation title from the first user message. It is
// best-effort: a failure to update the title never fails the message exchange.
func (s *Service) autoTitle(ctx context.Context, conv *domainai.Conversation, content string) {
	if len(conv.Messages) != 1 {
		return
	}
	title := content
	if len(title) > 50 {
		title = title[:50] + "..."
	}
	conv.Title = title
	_ = s.conversations.WithOrgConversations(ctx, func(conversations repository.ConversationRepository) error {
		return conversations.Update(ctx, conv)
	})
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
