CREATE TABLE IF NOT EXISTS friends (
    id BIGSERIAL PRIMARY KEY,
    sender_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ(0) NOT NULL DEFAULT NOW(),
    CONSTRAINT friends_status_valid CHECK (status IN ('pending', 'accepted')),
    CONSTRAINT friends_no_self CHECK (sender_id <> receiver_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_friends_unique_pair ON friends (LEAST(sender_id, receiver_id), GREATEST(sender_id, receiver_id));

CREATE INDEX IF NOT EXISTS idx_friends_sender_id ON friends (sender_id);
CREATE INDEX IF NOT EXISTS idx_friends_receiver_id ON friends (receiver_id);
