CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,
    receiver_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sender_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    post_id BIGINT REFERENCES posts(id) ON DELETE CASCADE,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ(0) NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_receiver_id ON notifications (receiver_id, created_at DESC);
CREATE INDEX idx_notifications_receiver_id_unread ON notifications (receiver_id) WHERE is_read = FALSE;
