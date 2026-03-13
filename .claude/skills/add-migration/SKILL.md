---
name: add-migration
description: Create numbered up/down SQL migration files following project conventions. Use when adding or modifying database tables.
user-invocable: true
allowed-tools: Read, Write, Bash, Glob, Grep
---

# Add Migration

Create migration for: **$ARGUMENTS**

You are creating a SQL migration for CleanSaaS. Follow every rule below precisely.

---

## STEP 1 — Determine the next migration number

Check existing migration files:
```bash
ls backend/migrations/*.up.sql 2>/dev/null | sort
```

If no migrations exist yet, start at `001`. Otherwise, increment the highest number by 1. Pad to 3 digits: `001`, `002`, ..., `010`, ..., `100`.

---

## STEP 2 — Determine the migration name

Parse `$ARGUMENTS` to derive a descriptive snake_case name.

Examples:
- "create users table" → `001_create_users`
- "add stripe_id to users" → `002_add_stripe_id_to_users`
- "create conversations table" → `003_create_conversations`
- "add index on email" → `004_add_index_on_users_email`

---

## STEP 3 — Create the UP migration

Create `backend/migrations/{NNN}_{name}.up.sql`

### SQL conventions (mandatory):

**Table creation:**
```sql
CREATE TABLE {table_name} (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- feature columns here
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

**Rules:**
- UUID primary keys with `gen_random_uuid()`
- `TEXT` not `VARCHAR` for string columns
- `TIMESTAMP NOT NULL DEFAULT NOW()` for created_at/updated_at
- Foreign key to users: `user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE`
- **NO cross-feature foreign keys** — only reference `users` table
- Index ALL foreign keys: `CREATE INDEX idx_{table}_{column} ON {table}({column});`
- Index frequently queried columns
- Use `NOT NULL` with sensible defaults where possible
- `BOOLEAN` columns default to `false`
- `JSONB` for flexible/nested data (e.g., features list, metadata)

**Column modifications:**
```sql
ALTER TABLE {table} ADD COLUMN {column} TEXT NOT NULL DEFAULT '';
```

**Index creation:**
```sql
CREATE INDEX idx_{table}_{column} ON {table}({column});
CREATE UNIQUE INDEX idx_{table}_{column}_unique ON {table}({column});
```

---

## STEP 4 — Create the DOWN migration

Create `backend/migrations/{NNN}_{name}.down.sql`

The down migration must **perfectly reverse** the up migration:

- `CREATE TABLE` → `DROP TABLE IF EXISTS {table};`
- `ALTER TABLE ADD COLUMN` → `ALTER TABLE {table} DROP COLUMN IF EXISTS {column};`
- `CREATE INDEX` → `DROP INDEX IF EXISTS {index_name};`
- `ALTER TABLE ADD CONSTRAINT` → `ALTER TABLE {table} DROP CONSTRAINT IF EXISTS {constraint};`

Always use `IF EXISTS` in down migrations to make them idempotent.

---

## STEP 5 — Verify the cross-feature FK rule

Read the up migration and check EVERY `REFERENCES` clause:
- `REFERENCES users(id)` → OK (users is core)
- `REFERENCES {anything_else}` → FORBIDDEN unless it's a self-reference or within the same feature

If a cross-feature FK is detected, restructure:
- Store the referenced ID as a plain column (no FK constraint)
- Add a comment explaining why there's no FK: `-- no FK: {feature} is independent`
- Or use an event/message pattern for cross-feature data needs

---

## STEP 6 — Feature prefix verification

Verify the table name follows feature conventions:
- Auth feature: `users`, `sessions`, `oauth_tokens`
- Billing feature: `subscriptions`, `plans`, `invoices`
- AI feature: `conversations`, `messages`
- Notification feature: `notifications`, `notification_templates`
- Storage feature: `files`
- Admin feature: `audit_logs`

New features should use a clear, descriptive table name that doesn't conflict with existing features.

---

## STEP 7 — Validate SQL syntax

Read the generated SQL and verify:
- No typos in SQL keywords
- Matching parentheses
- Correct PostgreSQL syntax (not MySQL/SQLite)
- Parameterized-ready schema (no hardcoded values that should be dynamic)
- Proper comma separation between columns

---

## Output

Report:
1. Files created: `backend/migrations/{NNN}_{name}.up.sql` and `.down.sql`
2. Tables/columns affected
3. Indexes created
4. FK verification: pass/fail
5. Any warnings or notes

Example:
```
Created:
  backend/migrations/003_create_conversations.up.sql
  backend/migrations/003_create_conversations.down.sql

Tables: conversations, messages
Indexes: idx_conversations_user_id, idx_messages_conversation_id
FK check: PASS (only references users table)
```
