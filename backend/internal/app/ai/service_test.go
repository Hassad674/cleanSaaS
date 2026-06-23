package ai

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	domainai "github.com/hassad/boilerplateSaaS/backend/internal/domain/ai"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/service"
)

// Mocks

type mockConvoRepo struct {
	createFn       func(ctx context.Context, c *domainai.Conversation) error
	findByIDFn     func(ctx context.Context, id string) (*domainai.Conversation, error)
	updateFn       func(ctx context.Context, c *domainai.Conversation) error
	deleteFn       func(ctx context.Context, id string) error
	listByUserIDFn func(ctx context.Context, userID string, offset, limit int) ([]*domainai.Conversation, int, error)
	addMessageFn   func(ctx context.Context, conversationID string, msg domainai.Message) error
}

func (m *mockConvoRepo) Create(ctx context.Context, c *domainai.Conversation) error {
	if m.createFn != nil {
		return m.createFn(ctx, c)
	}
	c.ID = "conv-1"
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}

func (m *mockConvoRepo) FindByID(ctx context.Context, id string) (*domainai.Conversation, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, domain.ErrNotFound
}

func (m *mockConvoRepo) Update(ctx context.Context, c *domainai.Conversation) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, c)
	}
	return nil
}

func (m *mockConvoRepo) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockConvoRepo) ListByUserID(ctx context.Context, userID string, offset, limit int) ([]*domainai.Conversation, int, error) {
	if m.listByUserIDFn != nil {
		return m.listByUserIDFn(ctx, userID, offset, limit)
	}
	return nil, 0, nil
}

func (m *mockConvoRepo) AddMessage(ctx context.Context, conversationID string, msg domainai.Message) error {
	if m.addMessageFn != nil {
		return m.addMessageFn(ctx, conversationID, msg)
	}
	return nil
}

// scopeOf adapts a plain mockConvoRepo into a repository.ConversationScope by
// invoking the callback with the underlying repo — unit tests need no real tx.
type mockConvoScope struct{ repo *mockConvoRepo }

func (s *mockConvoScope) WithOrgConversations(ctx context.Context, fn func(conversations repository.ConversationRepository) error) error {
	return fn(s.repo)
}

func scopeOf(repo *mockConvoRepo) *mockConvoScope { return &mockConvoScope{repo: repo} }

type mockAISvc struct {
	chatFn   func(ctx context.Context, messages []domainai.Message) (*service.AIResponse, error)
	streamFn func(ctx context.Context, messages []domainai.Message, writer io.Writer) error
}

func (m *mockAISvc) Chat(ctx context.Context, messages []domainai.Message) (*service.AIResponse, error) {
	if m.chatFn != nil {
		return m.chatFn(ctx, messages)
	}
	return &service.AIResponse{Content: "Hello!", Model: domainai.ModelGemini}, nil
}

func (m *mockAISvc) Stream(ctx context.Context, messages []domainai.Message, writer io.Writer) error {
	if m.streamFn != nil {
		return m.streamFn(ctx, messages, writer)
	}
	writer.Write([]byte("streamed response"))
	return nil
}

// Tests

func TestAIService_SendMessage_SavesToDB(t *testing.T) {
	var savedMessages []domainai.Message
	convoRepo := &mockConvoRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainai.Conversation, error) {
			return &domainai.Conversation{
				ID:     "conv-1",
				UserID: "user-1",
				Title:  "New conversation",
			}, nil
		},
		addMessageFn: func(_ context.Context, _ string, msg domainai.Message) error {
			savedMessages = append(savedMessages, msg)
			return nil
		},
	}
	aiMock := &mockAISvc{
		chatFn: func(_ context.Context, _ []domainai.Message) (*service.AIResponse, error) {
			return &service.AIResponse{Content: "AI response"}, nil
		},
	}

	svc := NewService(scopeOf(convoRepo), aiMock)
	reply, err := svc.SendMessage(context.Background(), "user-1", "conv-1", "Hello AI")

	assert.NoError(t, err)
	assert.Equal(t, "AI response", reply)
	assert.Len(t, savedMessages, 2)
	assert.Equal(t, domainai.RoleUser, savedMessages[0].Role)
	assert.Equal(t, "Hello AI", savedMessages[0].Content)
	assert.Equal(t, domainai.RoleAssistant, savedMessages[1].Role)
	assert.Equal(t, "AI response", savedMessages[1].Content)
}

func TestAIService_SendMessage_ForbiddenForOtherUser(t *testing.T) {
	convoRepo := &mockConvoRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainai.Conversation, error) {
			return &domainai.Conversation{
				ID:     "conv-1",
				UserID: "user-2",
			}, nil
		},
	}

	svc := NewService(scopeOf(convoRepo), &mockAISvc{})
	_, err := svc.SendMessage(context.Background(), "user-1", "conv-1", "Hello")

	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestAIService_DeleteConversation_OwnershipCheck(t *testing.T) {
	convoRepo := &mockConvoRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainai.Conversation, error) {
			return &domainai.Conversation{
				ID:     "conv-1",
				UserID: "user-2",
			}, nil
		},
	}

	svc := NewService(scopeOf(convoRepo), &mockAISvc{})
	err := svc.DeleteConversation(context.Background(), "user-1", "conv-1")

	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestAIService_CreateConversation_Success(t *testing.T) {
	convoRepo := &mockConvoRepo{}
	svc := NewService(scopeOf(convoRepo), &mockAISvc{})

	conv, err := svc.CreateConversation(context.Background(), "user-1", "My Chat")

	assert.NoError(t, err)
	assert.Equal(t, "My Chat", conv.Title)
	assert.Equal(t, "user-1", conv.UserID)
}
