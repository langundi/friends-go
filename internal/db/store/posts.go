package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Post struct {
	ID         int64
	UserID     int64
	Caption    string
	ImageURL   string
	ObjectKey  string
	LikeCount  int
	ReplyCount int
	CreatedAt  time.Time
	LikedByMe  bool
}

type PostStore struct {
	db *pgxpool.Pool
}

func NewPostStore(db *pgxpool.Pool) *PostStore {
	return &PostStore{db: db}
}

func (s *PostStore) CreatePost(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (user_id, caption, image_url, object_key)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := s.db.QueryRow(ctx, query,
		post.UserID,
		post.Caption,
		post.ImageURL,
		post.ObjectKey,
	).Scan(&post.ID, &post.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) DeletePostByID(ctx context.Context, id int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	_, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) GetTimeline(ctx context.Context, userID int64) ([]Post, error) {
	query := `
		SELECT p.id, p.user_id, p.caption, p.image_url, p.object_key, p.like_count, p.reply_count, p.created_at,
			EXISTS (
				SELECT 1 FROM likes l WHERE l.post_id = p.id AND l.user_id = $1
			) AS liked_by_me
		FROM posts p
		WHERE p.user_id = $1
			OR p.user_id IN (
				SELECT CASE
					WHEN f.sender_id = $1 THEN f.receiver_id
					ELSE f.sender_id
				END
				FROM friends f
				WHERE (f.sender_id = $1 or f.receiver_id = $1)
				AND f.status = 'accepted'
			)
		ORDER BY p.created_at DESC
		LIMIT 10
	`

	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	posts, err := pgx.CollectRows(rows, pgx.RowToStructByName[Post])
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *PostStore) GetPostByID(ctx context.Context, id int64) (*Post, error) {
	query := `
		SELECT p.id, p.user_id, p.caption, p.image_url, p.object_key, p.like_count, p.reply_count, p.created_at,
			EXISTS (
				SELECT 1 FROM likes l WHERE l.post_id = p.id AND l.user_id = $1
			) AS liked_by_me
		FROM posts p
		WHERE id = $1
	`

	var post Post

	err := s.db.QueryRow(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Caption,
		&post.ImageURL,
		&post.ObjectKey,
		&post.LikeCount,
		&post.ReplyCount,
		&post.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get post id %v:", id)
		}
		return nil, err
	}

	return &post, nil
}

func (s *PostStore) GetPostsByUserID(ctx context.Context, userID int64) ([]Post, error) {
	query := `
		SELECT p.id, p.user_id, p.caption, p.image_url, p.object_key, p.like_count, p.reply_count, p.created_at,
			EXISTS (
				SELECT 1 FROM likes l WHERE l.post_id = p.id AND l.user_id = $1
			) AS liked_by_me
		FROM posts p
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	posts, err := pgx.CollectRows(rows, pgx.RowToStructByName[Post])
	if err != nil {
		return nil, err
	}

	return posts, nil
}
