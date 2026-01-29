DROP INDEX IF EXISTS idx_contact_requests_deleted_at;
ALTER TABLE contact_requests DROP COLUMN IF EXISTS deleted_at;

