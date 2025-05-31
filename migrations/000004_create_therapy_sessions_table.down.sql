-- Remove the constraint on messages table first
ALTER TABLE messages DROP CONSTRAINT IF EXISTS fk_session_id;

-- Drop the therapy_sessions table
DROP TABLE IF EXISTS therapy_sessions;
