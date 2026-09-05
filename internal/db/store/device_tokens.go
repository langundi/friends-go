package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeviceToken struct {
	ID     int64
	UserID int64
	Token  string
}

type DeviceTokenStore struct {
	db *pgxpool.Pool
}

func NewDeviceTokenStore(db *pgxpool.Pool) *DeviceTokenStore {
	return &DeviceTokenStore{db: db}
}

func (s *DeviceTokenStore) CreateDeviceToken(ctx context.Context, userID int64, token string) error {
	query := `
		INSERT INTO device_tokens (user_id, token)
		VALUES ($1, $2)
	`
	_, err := s.db.Exec(ctx, query, userID, token)
	if err != nil {
		return fmt.Errorf("create device tokens: %w", err)
	}
	return nil
}

func (s *DeviceTokenStore) DeleteDeviceToken(ctx context.Context, userID int64, token string) error {
	query := `
		DELETE FROM device_tokens WHERE user_id = $1 AND token = $2
	`
	_, err := s.db.Exec(ctx, query, userID, token)
	if err != nil {
		return fmt.Errorf("delete device token: %w", err)
	}
	return nil
}

func (s *DeviceTokenStore) GetDeviceTokens(ctx context.Context, userID int64) ([]DeviceToken, error) {
	query := `
		SELECT id, user_id, token
		FROM device_tokens
		WHERE user_id = $1
	`

	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query device tokens: %w", err)
	}

	tokens, err := pgx.CollectRows(rows, pgx.RowToStructByName[DeviceToken])
	if err != nil {
		return nil, fmt.Errorf("collect device tokens: %w", err)
	}

	return tokens, nil
}
