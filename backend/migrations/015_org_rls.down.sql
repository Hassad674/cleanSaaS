-- Reverse migration 015: drop RLS policies, disable RLS, revoke grants, drop role.
DROP POLICY IF EXISTS org_isolation ON messages;
DROP POLICY IF EXISTS org_isolation ON notifications;
DROP POLICY IF EXISTS org_isolation ON conversations;
DROP POLICY IF EXISTS org_isolation ON files;
DROP POLICY IF EXISTS org_isolation ON subscriptions;

ALTER TABLE messages NO FORCE ROW LEVEL SECURITY;
ALTER TABLE messages DISABLE ROW LEVEL SECURITY;
ALTER TABLE notifications NO FORCE ROW LEVEL SECURITY;
ALTER TABLE notifications DISABLE ROW LEVEL SECURITY;
ALTER TABLE conversations NO FORCE ROW LEVEL SECURITY;
ALTER TABLE conversations DISABLE ROW LEVEL SECURITY;
ALTER TABLE files NO FORCE ROW LEVEL SECURITY;
ALTER TABLE files DISABLE ROW LEVEL SECURITY;
ALTER TABLE subscriptions NO FORCE ROW LEVEL SECURITY;
ALTER TABLE subscriptions DISABLE ROW LEVEL SECURITY;

-- Revoke privileges before dropping the role (a role owning grants cannot be
-- dropped). REVOKE is idempotent and safe even if some grants are absent.
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'app_user') THEN
        REVOKE ALL ON subscriptions, files, conversations, notifications, messages,
            organizations, organization_members, plans, users FROM app_user;
        REVOKE USAGE ON SCHEMA public FROM app_user;
        DROP ROLE app_user;
    END IF;
END $$;
