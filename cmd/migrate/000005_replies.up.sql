CREATE TABLE IF NOT EXISTS replies (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    reply TEXT NOT NULL,
    created_at TIMESTAMPTZ(0) NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_replies_post_id_created_at ON replies (post_id, created_at);
