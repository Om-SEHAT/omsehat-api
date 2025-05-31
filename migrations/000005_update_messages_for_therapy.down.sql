-- Remove the foreign key constraint for therapy sessions
ALTER TABLE messages DROP CONSTRAINT IF EXISTS fk_therapy_session_id;

-- Remove the index on the session_id column
DROP INDEX IF EXISTS idx_messages_session_id;

-- Remove the comment on the role column
COMMENT ON COLUMN messages.role IS NULL;

-- Restore the original size of the role column
ALTER TABLE messages ALTER COLUMN role TYPE VARCHAR(10);
