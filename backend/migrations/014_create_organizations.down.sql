-- Reverse migration 014: drop org_id from tenant tables, then drop org tables.
ALTER TABLE subscriptions DROP COLUMN IF EXISTS org_id;
ALTER TABLE files         DROP COLUMN IF EXISTS org_id;
ALTER TABLE conversations DROP COLUMN IF EXISTS org_id;
ALTER TABLE notifications DROP COLUMN IF EXISTS org_id;

DROP TABLE IF EXISTS organization_members;
DROP TABLE IF EXISTS organizations;
