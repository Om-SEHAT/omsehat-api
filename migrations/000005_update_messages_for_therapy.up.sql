-- Add a new type for the role in messages table to support therapist role
ALTER TABLE messages ALTER COLUMN role TYPE VARCHAR(50);

-- Add a comment to explain the role column values
COMMENT ON COLUMN messages.role IS 'Role can be: user, omsapa, or therapist';

-- Add a foreign key constraint for therapy sessions
ALTER TABLE messages 
ADD CONSTRAINT fk_therapy_session_id 
FOREIGN KEY (session_id) 
REFERENCES therapy_sessions(id) 
ON DELETE CASCADE;

-- Create an index on the session_id column for faster lookups
CREATE INDEX idx_messages_session_id ON messages(session_id);
