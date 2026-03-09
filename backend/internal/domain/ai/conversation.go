package ai

import "time"

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
)

type Message struct {
	Role    Role
	Content string
}

type Conversation struct {
	ID        string
	UserID    string
	Title     string
	Model     Model
	Messages  []Message
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Conversation) AddMessage(role Role, content string) {
	c.Messages = append(c.Messages, Message{Role: role, Content: content})
	c.UpdatedAt = time.Now()
}
