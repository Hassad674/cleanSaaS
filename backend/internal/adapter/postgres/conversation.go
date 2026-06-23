package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hassad/boilerplateSaaS/backend/internal/domain"
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/ai"
	"github.com/hassad/boilerplateSaaS/backend/pkg/orgctx"
)

// ConversationRepository implements repository.ConversationRepository. It holds a
// DBTX so the same code runs on the pool or an org-scoped transaction. Every query
// filters by the active org_id (defense layer 2) on top of RLS (layer 3); inserts
// stamp the active org_id onto the row.
type ConversationRepository struct {
	db DBTX
}

func NewConversationRepository(db *sql.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

// newConversationRepositoryTx binds the repository to an open transaction (org scope).
func newConversationRepositoryTx(tx DBTX) *ConversationRepository {
	return &ConversationRepository{db: tx}
}

func (r *ConversationRepository) Create(ctx context.Context, c *ai.Conversation) error {
	orgID, ok := orgctx.OrgID(ctx)
	if !ok {
		return fmt.Errorf("creating conversation: %w", domain.ErrForbidden)
	}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO conversations (user_id, org_id, title) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`,
		c.UserID, orgID, c.Title,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return fmt.Errorf("creating conversation: %w", err)
	}
	return nil
}

func (r *ConversationRepository) FindByID(ctx context.Context, id string) (*ai.Conversation, error) {
	c := &ai.Conversation{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, title, created_at, updated_at FROM conversations WHERE id = $1 AND org_id = $2`, id, orgFilter(ctx),
	).Scan(&c.ID, &c.UserID, &c.Title, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("finding conversation: %w", err)
	}

	// Load messages. Messages have no org_id of their own; they are scoped
	// transitively through their conversation (which we just confirmed is visible
	// under the active org), and PostgreSQL RLS enforces the same rule on the wire.
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
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating messages: %w", err)
	}

	return c, nil
}

func (r *ConversationRepository) Update(ctx context.Context, c *ai.Conversation) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE conversations SET title = $1, updated_at = NOW() WHERE id = $2 AND org_id = $3`,
		c.Title, c.ID, orgFilter(ctx),
	)
	if err != nil {
		return fmt.Errorf("updating conversation: %w", err)
	}
	return nil
}

func (r *ConversationRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM conversations WHERE id = $1 AND org_id = $2`, id, orgFilter(ctx))
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
	org := orgFilter(ctx)
	var total int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM conversations WHERE user_id = $1 AND org_id = $2`, userID, org,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("counting conversations: %w", err)
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, title, created_at, updated_at
		 FROM conversations WHERE user_id = $1 AND org_id = $2
		 ORDER BY updated_at DESC
		 LIMIT $3 OFFSET $4`, userID, org, limit, offset,
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
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterating conversations: %w", err)
	}

	return convos, total, nil
}

func (r *ConversationRepository) AddMessage(ctx context.Context, conversationID string, msg ai.Message) error {
	// Confirm the conversation belongs to the active org before writing a message.
	// RLS also enforces this transitively, but this gives a clear domain error
	// instead of a silently-rejected insert when the conversation is cross-tenant.
	var exists bool
	if err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM conversations WHERE id = $1 AND org_id = $2)`,
		conversationID, orgFilter(ctx),
	).Scan(&exists); err != nil {
		return fmt.Errorf("verifying conversation org: %w", err)
	}
	if !exists {
		return domain.ErrNotFound
	}

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO messages (conversation_id, role, content) VALUES ($1, $2, $3)`,
		conversationID, string(msg.Role), msg.Content,
	)
	if err != nil {
		return fmt.Errorf("adding message: %w", err)
	}
	// Update conversation timestamp (org-scoped).
	if _, err := r.db.ExecContext(ctx,
		`UPDATE conversations SET updated_at = NOW() WHERE id = $1 AND org_id = $2`, conversationID, orgFilter(ctx),
	); err != nil {
		return fmt.Errorf("updating conversation timestamp: %w", err)
	}
	return nil
}
