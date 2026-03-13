package ai

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConversation_AddMessage(t *testing.T) {
	conv := &Conversation{
		ID:        "conv-1",
		UserID:    "user-1",
		Title:     "Test Chat",
		Model:     ModelClaude,
		Messages:  nil,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now().Add(-time.Hour),
	}

	before := conv.UpdatedAt
	conv.AddMessage(RoleUser, "Hello!")

	assert.Len(t, conv.Messages, 1)
	assert.Equal(t, RoleUser, conv.Messages[0].Role)
	assert.Equal(t, "Hello!", conv.Messages[0].Content)
	assert.True(t, conv.UpdatedAt.After(before), "UpdatedAt should advance after AddMessage")
}

func TestConversation_AddMessage_MultipleMessages(t *testing.T) {
	conv := &Conversation{
		ID:     "conv-1",
		UserID: "user-1",
	}

	conv.AddMessage(RoleUser, "Hi")
	conv.AddMessage(RoleAssistant, "Hello!")
	conv.AddMessage(RoleUser, "How are you?")

	assert.Len(t, conv.Messages, 3)
	assert.Equal(t, RoleUser, conv.Messages[0].Role)
	assert.Equal(t, RoleAssistant, conv.Messages[1].Role)
	assert.Equal(t, RoleUser, conv.Messages[2].Role)
}

func TestConversation_AddMessage_SystemRole(t *testing.T) {
	conv := &Conversation{
		ID:     "conv-1",
		UserID: "user-1",
	}

	conv.AddMessage(RoleSystem, "You are a helpful assistant.")
	assert.Len(t, conv.Messages, 1)
	assert.Equal(t, RoleSystem, conv.Messages[0].Role)
	assert.Equal(t, "You are a helpful assistant.", conv.Messages[0].Content)
}

func TestMessage_Fields(t *testing.T) {
	msg := Message{
		Role:    RoleUser,
		Content: "test content",
	}

	assert.Equal(t, RoleUser, msg.Role)
	assert.Equal(t, "test content", msg.Content)
}

func TestRole_Constants(t *testing.T) {
	assert.Equal(t, Role("user"), RoleUser)
	assert.Equal(t, Role("assistant"), RoleAssistant)
	assert.Equal(t, Role("system"), RoleSystem)
}
