package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/ai"
)

type ConversationRepository struct {
	db *sql.DB
}

func NewConversationRepository(db *sql.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

func (r *ConversationRepository) Create(ctx context.Context, c *ai.Conversation) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO conversations (user_id, title) VALUES ($1, $2) RETURNING id, created_at, updated_at`,
		c.UserID, c.Title,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return fmt.Errorf("creating conversation: %w", err)
	}
	return nil
}

func (r *ConversationRepository) FindByID(ctx context.Context, id string) (*ai.Conversation, error) {
	c := &ai.Conversation{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, title, created_at, updated_at FROM conversations WHERE id = $1`, id,
	).Scan(&c.ID, &c.UserID, &c.Title, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("finding conversation: %w", err)
	}

	// Load messages
	rows, err := r.db.QueryContext(ctx,
		`SELECT role, content FROM messages WHERE conversation_id = $1 ORDER BY created_at ASC`, id,
	)
	if err != nil {
		return nil, fmt.Errorf("loading messages: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var m ai.Message
		var role string
		if err := rows.Scan(&role, &m.Content); err != nil {
			return nil, fmt.Errorf("scanning message: %w", err)
		}
		m.Role = ai.Role(role)
		c.Messages = append(c.Messages, m)
	}

	return c, nil
}

func (r *ConversationRepository) Update(ctx context.Context, c *ai.Conversation) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE conversations SET title = $1, updated_at = NOW() WHERE id = $2`,
		c.Title, c.ID,
	)
	if err != nil {
		return fmt.Errorf("updating conversation: %w", err)
	}
	return nil
}

func (r *ConversationRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM conversations WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("deleting conversation: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ConversationRepository) ListByUserID(ctx context.Context, userID string, offset, limit int) ([]*ai.Conversation, int, error) {
	var total int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM conversations WHERE user_id = $1`, userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("counting conversations: %w", err)
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, title, created_at, updated_at
		 FROM conversations WHERE user_id = $1
		 ORDER BY updated_at DESC
		 LIMIT $2 OFFSET $3`, userID, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("listing conversations: %w", err)
	}
	defer rows.Close()

	var convos []*ai.Conversation
	for rows.Next() {
		c := &ai.Conversation{}
		if err := rows.Scan(&c.ID, &c.UserID, &c.Title, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scanning conversation: %w", err)
		}
		convos = append(convos, c)
	}

	return convos, total, nil
}

func (r *ConversationRepository) AddMessage(ctx context.Context, conversationID string, msg ai.Message) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO messages (conversation_id, role, content) VALUES ($1, $2, $3)`,
		conversationID, string(msg.Role), msg.Content,
	)
	if err != nil {
		return fmt.Errorf("adding message: %w", err)
	}
	// Update conversation timestamp
	_, _ = r.db.ExecContext(ctx,
		`UPDATE conversations SET updated_at = NOW() WHERE id = $1`, conversationID,
	)
	return nil
}
