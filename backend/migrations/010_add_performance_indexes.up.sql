-- Performance indexes for columns used in WHERE, ORDER BY, and JOIN clauses
-- that were missing from the original migrations.

-- blog_posts: author_id is a FK but was never indexed (needed for joins and author lookups)
CREATE INDEX IF NOT EXISTS idx_blog_posts_author_id ON blog_posts(author_id);

-- blog_posts: GIN index on tags for efficient ANY(tags) queries
CREATE INDEX IF NOT EXISTS idx_blog_posts_tags ON blog_posts USING GIN(tags);

-- conversations: updated_at used in ORDER BY for ListByUserID
CREATE INDEX IF NOT EXISTS idx_conversations_user_updated ON conversations(user_id, updated_at DESC);

-- notifications: created_at used in ORDER BY for list queries
CREATE INDEX IF NOT EXISTS idx_notifications_user_created ON notifications(user_id, created_at DESC);

-- invoices: created_at used in ORDER BY for list queries
CREATE INDEX IF NOT EXISTS idx_invoices_user_created ON invoices(user_id, created_at DESC);

-- files: created_at used in ORDER BY for list queries
CREATE INDEX IF NOT EXISTS idx_files_user_created ON files(user_id, created_at DESC);

-- messages: created_at used in ORDER BY for loading messages
CREATE INDEX IF NOT EXISTS idx_messages_conversation_created ON messages(conversation_id, created_at ASC);

-- subscriptions: composite index for common lookup pattern (user's active subscription)
CREATE INDEX IF NOT EXISTS idx_subscriptions_user_status ON subscriptions(user_id, status);

-- password_resets: expires_at used in cleanup queries
CREATE INDEX IF NOT EXISTS idx_password_resets_expires_at ON password_resets(expires_at);

-- email_verifications: expires_at used in cleanup queries
CREATE INDEX IF NOT EXISTS idx_email_verifications_expires_at ON email_verifications(expires_at);
