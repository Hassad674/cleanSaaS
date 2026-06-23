-- Row-Level Security: the LAST line of org-isolation defense (layer 3).
--
-- Even if a repository forgot its `WHERE org_id = $n` filter, these policies make
-- the database itself refuse to return or write another tenant's rows. The policy
-- is keyed on a per-transaction GUC `app.current_org_id` set by the request path
-- (SET LOCAL app.current_org_id = '<uuid>') together with SET LOCAL ROLE app_user.
--
-- DENY BY DEFAULT: current_setting('app.current_org_id', true) returns NULL when
-- the GUC was never set in the session (the `true` = missing_ok argument). Once a
-- connection has set it at least once it instead reads back as '' after reset, so
-- NULLIF(..., '') folds both "never set" and "reset to empty" to NULL. `org_id =
-- NULL` is NULL (never true), so a missing org context matches NO rows and rejects
-- every write — it leaks nothing and, thanks to NULLIF, never errors on an empty
-- GUC either.
--
-- The app_user role is created NOBYPASSRLS (and is not a superuser/table owner),
-- so RLS is actually enforced for it. System paths (migrations, seed, scheduler
-- cleanups) keep using the privileged connection role, which owns the tables and
-- therefore bypasses RLS — exactly what those trusted, cross-tenant jobs need.

-- 1. Restricted application role. Created idempotently; NOBYPASSRLS is the point.
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'app_user') THEN
        CREATE ROLE app_user NOLOGIN NOBYPASSRLS;
    ELSE
        ALTER ROLE app_user NOBYPASSRLS;
    END IF;
END $$;

-- 2. Privileges. app_user needs to use the schema and read/write the tenant tables
-- (and their child tables). It does NOT own them, so RLS applies.
GRANT USAGE ON SCHEMA public TO app_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON
    subscriptions, files, conversations, notifications, messages,
    organizations, organization_members
TO app_user;
-- Read-only access to lookup tables the tenant request path needs (plans is a
-- global catalog; users is needed for joins). No write access is granted.
GRANT SELECT ON plans, users TO app_user;

-- 3. Enable + FORCE RLS on the tenant-scoped tables. FORCE makes the policy apply
-- even to the table owner when it connects as a non-superuser, closing the
-- "owner silently bypasses" gap for any future owner-role connection.
ALTER TABLE subscriptions ENABLE ROW LEVEL SECURITY;
ALTER TABLE subscriptions FORCE ROW LEVEL SECURITY;
ALTER TABLE files ENABLE ROW LEVEL SECURITY;
ALTER TABLE files FORCE ROW LEVEL SECURITY;
ALTER TABLE conversations ENABLE ROW LEVEL SECURITY;
ALTER TABLE conversations FORCE ROW LEVEL SECURITY;
ALTER TABLE notifications ENABLE ROW LEVEL SECURITY;
ALTER TABLE notifications FORCE ROW LEVEL SECURITY;

-- 4. Deny-by-default org-isolation policies. USING governs visibility (SELECT /
-- UPDATE / DELETE row access); WITH CHECK governs new/modified row values
-- (INSERT / UPDATE), so a tenant can neither read nor write across the boundary.
CREATE POLICY org_isolation ON subscriptions
    USING (org_id = NULLIF(current_setting('app.current_org_id', true), '')::uuid)
    WITH CHECK (org_id = NULLIF(current_setting('app.current_org_id', true), '')::uuid);
CREATE POLICY org_isolation ON files
    USING (org_id = NULLIF(current_setting('app.current_org_id', true), '')::uuid)
    WITH CHECK (org_id = NULLIF(current_setting('app.current_org_id', true), '')::uuid);
CREATE POLICY org_isolation ON conversations
    USING (org_id = NULLIF(current_setting('app.current_org_id', true), '')::uuid)
    WITH CHECK (org_id = NULLIF(current_setting('app.current_org_id', true), '')::uuid);
CREATE POLICY org_isolation ON notifications
    USING (org_id = NULLIF(current_setting('app.current_org_id', true), '')::uuid)
    WITH CHECK (org_id = NULLIF(current_setting('app.current_org_id', true), '')::uuid);

-- 5. messages have no org_id of their own; they belong to their parent
-- conversation. Scope them transitively: a message is visible/insertable only if
-- its conversation is visible under the active org. Because the EXISTS subquery
-- reads conversations (itself RLS-protected), cross-tenant message access is
-- rejected without duplicating org_id onto every message row.
ALTER TABLE messages ENABLE ROW LEVEL SECURITY;
ALTER TABLE messages FORCE ROW LEVEL SECURITY;
CREATE POLICY org_isolation ON messages
    USING (EXISTS (SELECT 1 FROM conversations c WHERE c.id = messages.conversation_id))
    WITH CHECK (EXISTS (SELECT 1 FROM conversations c WHERE c.id = messages.conversation_id));
