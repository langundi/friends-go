package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshToken struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	Revoked   bool      `json:"revoked"`
}

type RefeshTokenStore struct {
	db *pgxpool.Pool
}

func NewRefreshTokenStore(db *pgxpool.Pool) *RefeshTokenStore {
	return &RefeshTokenStore{db: db}
}

func (s *RefeshTokenStore) CreateRefreshToken(ctx context.Context, token *RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at, revoked)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	return s.db.QueryRow(ctx, query,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.Revoked,
	).Scan(&token.ID, &token.CreatedAt)
}

func (s *RefeshTokenStore) GetRefreshToken(ctx context.Context, tokenString string) (*RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at, revoked
		FROM refresh_tokens
		WHERE token = $1
	`

	var token RefreshToken

	err := s.db.QueryRow(ctx, query, tokenString).Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.Revoked,
	)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (s *RefeshTokenStore) DeleteRefreshToken(ctx context.Context, tokenString string) error {
	query := `DELETE FROM refresh_token WHERE token = $1`

	_, err := s.db.Exec(ctx, query, tokenString)
	if err != nil {
		return err
	}

	return nil
}

func (s *RefeshTokenStore) RevokeRefreshToken(ctx context.Context, tokenString string) error {
	query := `
		UDPATE refresh_tokens
		SET revoked = true
		WHERE token = $1
	`

	_, err := s.db.Exec(ctx, query, tokenString)
	if err != nil {
		return err
	}

	return nil
}
