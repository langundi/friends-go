package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Like struct {
	UserID int64
	PostID int64
}

type LikeStore struct {
	db *pgxpool.Pool
}

func NewLikeStore(db *pgxpool.Pool) *LikeStore {
	return &LikeStore{db: db}
}

func (s *LikeStore) LikePost(ctx context.Context, userID, postID int64) error {
	query := `
		INSERT INTO likes (user_id, post_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, post_id) DO NOTHING
	`

	_, err := s.db.Exec(ctx, query, userID, postID)
	if err != nil {
		return err
	}

	return nil
}

func (s *LikeStore) UnlikePost(ctx context.Context, userID, postID int64) error {
	query := `
		DELETE FROM likes WHERE user_id = $1 AND post_id = $2
	`

	_, err := s.db.Exec(ctx, query, userID, postID)
	if err != nil {
		return err
	}

	return nil
}
