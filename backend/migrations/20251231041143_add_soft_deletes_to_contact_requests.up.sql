-- Add soft delete column to contact_requests table
ALTER TABLE contact_requests ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

-- Create index on deleted_at for better query performance
CREATE INDEX IF NOT EXISTS idx_contact_requests_deleted_at ON contact_requests (deleted_at) WHERE deleted_at IS NULL;

