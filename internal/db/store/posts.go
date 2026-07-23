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
	ID        int64
	UserID    int64
	Caption   string
	ImageURL  string
	CreatedAt time.Time
}

type PostStore struct {
	db *pgxpool.Pool
}

func NewPostStore(db *pgxpool.Pool) *PostStore {
	return &PostStore{db: db}
}

func (s *PostStore) CreatePost(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (user_id, caption, image_url)
		VALUES($1, $2, $3)
		RETURNING id, created_at
	`

	err := s.db.QueryRow(ctx, query,
		post.UserID,
		post.Caption,
		post.ImageURL,
	).Scan(&post.ID, &post.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) GetPostByID(ctx context.Context, id int64) (*Post, error) {
	query := `
		SELECT id, user_id, caption, image_url, created_at
		FROM posts
		WHERE id = $1
	`

	var post Post

	err := s.db.QueryRow(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Caption,
		&post.ImageURL,
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
