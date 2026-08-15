CREATE TABLE IF NOT EXISTS friends (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    friend_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ(0) NOT NULL DEFAULT NOW(),
    CONSTRAINT friends_status_valid CHECK (status IN ('pending', 'accepted')),
    CONSTRAINT friends_no_self CHECK (user_id <> friend_id),
    UNIQUE (user_id, friend_id)
);

CREATE INDEX IF NOT EXISTS idx_friends_friend_id ON friends (friend_id);
