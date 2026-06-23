-- Org-based multi-tenancy: an organization is the tenant. Every tenant-scoped
-- resource (subscriptions, files, conversations, notifications) belongs to one
-- organization. This migration creates the org tables and adds + indexes an
-- org_id column on the tenant-owned tables. RLS enforcement is layered on top in
-- migration 015. (blog stays GLOBAL/public; teams are left as-is for now.)

CREATE TABLE IF NOT EXISTS organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_organizations_owner_id ON organizations(owner_id);
CREATE INDEX IF NOT EXISTS idx_organizations_slug ON organizations(slug);

CREATE TABLE IF NOT EXISTS organization_members (
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role TEXT NOT NULL DEFAULT 'member',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (org_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_organization_members_user_id ON organization_members(user_id);

-- Add org_id to the tenant-owned tables. ON DELETE CASCADE: deleting an
-- organization removes its tenant data. The column is added nullable first so the
-- migration is safe to run against a populated DB, then constrained NOT NULL.
-- A boilerplate ships with a fresh DB (backfill is N/A), but a future deployment
-- that already has rows must backfill org_id BEFORE applying migration 015's
-- NOT NULL is harmless here because we set NOT NULL only when the table is empty
-- of unscoped rows; see the guard below.

ALTER TABLE subscriptions ADD COLUMN IF NOT EXISTS org_id UUID REFERENCES organizations(id) ON DELETE CASCADE;
ALTER TABLE files         ADD COLUMN IF NOT EXISTS org_id UUID REFERENCES organizations(id) ON DELETE CASCADE;
ALTER TABLE conversations ADD COLUMN IF NOT EXISTS org_id UUID REFERENCES organizations(id) ON DELETE CASCADE;
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS org_id UUID REFERENCES organizations(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_subscriptions_org_id ON subscriptions(org_id);
CREATE INDEX IF NOT EXISTS idx_files_org_id ON files(org_id);
CREATE INDEX IF NOT EXISTS idx_conversations_org_id ON conversations(org_id);
CREATE INDEX IF NOT EXISTS idx_notifications_org_id ON notifications(org_id);

-- Enforce NOT NULL only when no pre-existing unscoped rows would be violated.
-- On a fresh boilerplate DB the tables are empty, so this sets NOT NULL cleanly.
-- On a DB with legacy rows lacking org_id, the constraint is skipped (logged via
-- a NOTICE) so the migration never fails mid-flight; backfill then re-run.
DO $$
DECLARE
    t TEXT;
    unscoped BIGINT;
BEGIN
    FOREACH t IN ARRAY ARRAY['subscriptions', 'files', 'conversations', 'notifications'] LOOP
        EXECUTE format('SELECT count(*) FROM %I WHERE org_id IS NULL', t) INTO unscoped;
        IF unscoped = 0 THEN
            EXECUTE format('ALTER TABLE %I ALTER COLUMN org_id SET NOT NULL', t);
        ELSE
            RAISE NOTICE 'table % has % rows with NULL org_id; leaving org_id nullable — backfill then SET NOT NULL', t, unscoped;
        END IF;
    END LOOP;
END $$;
