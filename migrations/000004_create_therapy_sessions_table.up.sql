CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS therapy_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    stress_level INT NOT NULL,
    mood_rating INT NOT NULL,
    sleep_quality INT NOT NULL,
    is_health_worker BOOLEAN NOT NULL,
    specialization VARCHAR(100),
    work_hours INT,
    recommendation TEXT,
    next_steps TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Update the messages table to be compatible with therapy sessions
ALTER TABLE messages
    ADD CONSTRAINT fk_session_id
    FOREIGN KEY (session_id) 
    REFERENCES sessions(id)
    ON DELETE CASCADE;

-- Add a comment column to explain when a role is 'therapist'
COMMENT ON COLUMN messages.role IS 'Can be user, omsapa, or therapist';
