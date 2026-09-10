package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Reply struct {
	ID             int64
	UserID         int64
	PostID         int64
	Reply          string
	CreatedAt      time.Time
	RepliedByMe    bool
	Username       string
	ProfilePicture *string
}

type ReplyStore struct {
	db *pgxpool.Pool
}

func NewReplyStore(db *pgxpool.Pool) *ReplyStore {
	return &ReplyStore{db: db}
}

func (s *ReplyStore) CreateReply(ctx context.Context, reply *Reply) error {
	query := `
		INSERT INTO replies (user_id, post_id, reply)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := s.db.QueryRow(ctx, query,
		reply.UserID,
		reply.PostID,
		reply.Reply,
	).Scan(&reply.ID, &reply.CreatedAt)
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
		SELECT r.id, r.user_id, r.post_id, r.reply, r.created_at, r.user_id = $1 AS replied_by_me, u.username, u.profile_picture
		FROM replies r
		JOIN users u ON u.id = r.user_id
		WHERE post_id = $2
		ORDER BY created_at ASC
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
