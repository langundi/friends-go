package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Reply struct {
	ID          int64
	UserID      int64
	PostID      int64
	Reply       string
	CreatedAt   time.Time
	RepliedByMe bool
}

type ReplyStore struct {
	db *pgxpool.Pool
}

func NewReplyStore(db *pgxpool.Pool) *ReplyStore {
	return &ReplyStore{db: db}
}

func (s *ReplyStore) CreateReply(ctx context.Context, userID, postID int64, reply string) error {
	query := `
		INSERT INTO replies (user_id, post_id, reply)
		VALUES ($1, $2, $3)
	`

	_, err := s.db.Exec(ctx, query, userID, postID, reply)
	if err != nil {
		return err
	}

	return nil
}

func (s *ReplyStore) DeleteReply(ctx context.Context, replyID int64) error {
	query := `
		DELETE FROM replies WHERE id = $1
	`

	_, err := s.db.Exec(ctx, query, replyID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ReplyStore) GetRepliesForPostID(ctx context.Context, userID, postID int64) ([]Reply, error) {
	query := `
		SELECT id, user_id, post_id, reply, created_at, user_id = $1 AS replied_by_me
		FROM replies
		WHERE post_id = $2
		ORDER BY created_at ASC
		LIMIT 10
	`

	rows, err := s.db.Query(ctx, query, userID, postID)
	if err != nil {
		return nil, err
	}

	replies, err := pgx.CollectRows(rows, pgx.RowToStructByName[Reply])
	if err != nil {
		return nil, err
	}

	return replies, nil
}
