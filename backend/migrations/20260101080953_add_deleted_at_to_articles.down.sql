DROP INDEX IF EXISTS idx_articles_deleted_at;

ALTER TABLE articles DROP COLUMN IF EXISTS deleted_at;

