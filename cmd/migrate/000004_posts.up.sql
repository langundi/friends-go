CREATE TABLE IF NOT EXISTS posts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    caption TEXT,
    image_url TEXT NOT NULL,
    object_key TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ(0) NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_posts_created_at on posts (created_at DESC);
